package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/db/generated"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/domain"
	"github.com/rasm445f/soft-exam-2/handlers"
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
	handler := cors.Default().Handler(mux)

	return handler, err
}

func main() {
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
