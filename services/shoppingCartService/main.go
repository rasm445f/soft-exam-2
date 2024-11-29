package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/domain"
	"github.com/rasm445f/soft-exam-2/handlers"
)

func run() (http.Handler, error) {
	redisClient, err := db.Redis_conn()
	if err != nil {
		return nil, err
	}

	repo := db.NewShoppingCartRepository(redisClient)
	shoppingDomain := domain.NewShoppingCartDomain(repo)
	shoppingHandler := handlers.NewShoppingCartHandler(shoppingDomain)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /api/shopping", shoppingHandler.AddItem())
	mux.HandleFunc("PATCH /api/shopping/{customerId}/{itemId}", shoppingHandler.UpdateCart())
	mux.HandleFunc("GET /api/shopping/{customerId}", shoppingHandler.ViewCart())
	mux.HandleFunc("DELETE /api/shopping/{customerId}", shoppingHandler.ClearCart())
	// Broker
	mux.HandleFunc("GET /api/shopping/consume", shoppingHandler.ConsumeMenuItem())
	mux.HandleFunc("POST /api/shopping/publish/{customerId}", shoppingHandler.PublishShoppingCart())

	//CORS stuff
	handler := cors.Default().Handler(mux)

	return handler, nil
}

func main() {
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
