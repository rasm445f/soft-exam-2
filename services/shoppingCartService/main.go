package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/domain"
	"github.com/rasm445f/soft-exam-2/handlers"
	"github.com/rasm445f/soft-exam-2/metrics"
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
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /api/shopping", shoppingHandler.AddItem())
	mux.HandleFunc("PATCH /api/shopping/{customerId}/{itemId}", shoppingHandler.UpdateCart())
	mux.HandleFunc("GET /api/shopping/{customerId}", shoppingHandler.ViewCart())
	mux.HandleFunc("DELETE /api/shopping/{customerId}", shoppingHandler.ClearCart())
	// Broker
	mux.HandleFunc("GET /api/shopping/consume", shoppingHandler.ConsumeMenuItem())
	mux.HandleFunc("POST /api/shopping/publish/{customerId}", shoppingHandler.PublishShoppingCart())

	//CORS stuff
	metrics := metrics.MetricsMiddleware(mux)
	handler := cors.Default().Handler(metrics)

	//test change for cicd
	return handler, nil
}

// @title Order Shopping cart API
// @version 1.0
// @description This is the API documentation for the Shopping cart Service.
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8084
func main() {
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8084")
	log.Fatal(http.ListenAndServe(":8084", mux))
}
