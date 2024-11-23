package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/oTuff/go-startkode/db"
	_ "github.com/oTuff/go-startkode/docs"
	"github.com/oTuff/go-startkode/domain"
	"github.com/oTuff/go-startkode/handlers"
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
	// mux.HandleFunc("DEL /api/shopping/{customerId}", shoppingHandler.ClearCart())

	//CORS stuff
	handler := cors.Default().Handler(mux)

	return handler, nil
}

func main() {
	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
