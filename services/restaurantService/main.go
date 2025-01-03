package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/db/generated"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/domain"
	"github.com/rasm445f/soft-exam-2/handlers"
	"github.com/rasm445f/soft-exam-2/metrics"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func run() (http.Handler, error) {
	db, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Initialize queries and domain layer
	queries := generated.New(db)
	restaurantDomain := domain.NewRestaurantDomain(queries)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantDomain)

	mux := http.NewServeMux()

	// Routes
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/restaurants", restaurantHandler.GetAllRestaurants())
	mux.HandleFunc("GET /api/restaurants/{restaurantId}", restaurantHandler.GetRestaurantById())
	mux.HandleFunc("GET /api/restaurants/{restaurantId}/menu-items", restaurantHandler.GetMenuItemsByRestaurant())
	mux.HandleFunc("GET /api/restaurants/{restaurantId}/menu-items/{menuitemId}", restaurantHandler.GetMenuItemByRestaurantAndId())
	mux.HandleFunc("GET /api/categories", restaurantHandler.GetAllCategories())
	mux.HandleFunc("GET /api/filter/{category}", restaurantHandler.FilterRestaurantByCategory())
	// Broker
	mux.HandleFunc("POST /api/restaurants/menu/select", restaurantHandler.SelectMenuItem())

	//CORS stuff
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	metrics := metrics.MetricsMiddleware(mux)
	handler := corsHandler.Handler(metrics)

	return handler, err
}

// @title Restaurant Service API
// @version 1.0
// @description This is the API documentation for the Restaurant Service.
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8083
func main() {
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}
