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

	// Initialize Queries with DB
	queries := generated.New(db)
	orderDomain := domain.NewOrderDomain(queries)
	orderHandler := handlers.NewOrderHandler(orderDomain)
	feedbackDomain := domain.NewFeedbackDomain(queries)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackDomain)
	deliveryAgentDomain := domain.NewDeliveryAgentDomain(queries)
	deliveryAgentHandler := handlers.NewDeliveryAgentHandler(deliveryAgentDomain)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/orders", orderHandler.GetAllOrders())
	mux.HandleFunc("GET /api/orders/{orderId}", orderHandler.GetOrderById())
	mux.HandleFunc("PATCH /api/order/status/{orderId}", orderHandler.UpdateOrderStatus())
	mux.HandleFunc("DELETE /api/orders/{orderId}", orderHandler.DeleteOrder())
	mux.HandleFunc("PATCH /api/order/status-agent/{orderId}", orderHandler.UpdateOrderStatusWithDeliveryAgentId())
	mux.HandleFunc("GET /api/order/bonus/{orderId}", orderHandler.CalculateOrderBonus())
	// Feedback
	mux.HandleFunc("GET /api/feedbacks", feedbackHandler.GetAllFeedbacks())
	mux.HandleFunc("GET /api/feedbacks/{orderId}", feedbackHandler.GetFeedbackByOrderId())
	mux.HandleFunc("POST /api/feedback", feedbackHandler.CreateFeedback())
	// DeliveryAgent
	mux.HandleFunc("GET /api/delivery-agent", deliveryAgentHandler.GetAllDeliveryAgents())
	mux.HandleFunc("GET /api/delivery-agent/{deliveryAgentId}", deliveryAgentHandler.GetDeliveryAgentById())
	mux.HandleFunc("POST /api/delivery-agent", deliveryAgentHandler.CreateDeliveryAgent())
	// Broker
	mux.HandleFunc("GET /api/order/consume", orderHandler.ConsumeOrder())

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

	fmt.Println("Running server on port 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
