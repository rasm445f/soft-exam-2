package handlers

//TODO: implement system test
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type Customer struct {
	ID          int32   `json:"id"`
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	Phonenumber *string `json:"phonenumber"`
	Addressid   *int32  `json:"addressid"`
}

func TestSystem(t *testing.T) {
	// make sure all the services are running first
	client := &http.Client{}

	// create a customer
	customerPayload := `{
		"email": "test@test.dk",
		"name": "string",
		"password": "Test123!!!!",
		"phonenumber": "12345678",
		"street_address": "ligma",
		"zip_code": 2800
	}`

	// // Create the request
	req1, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/customer", bytes.NewBuffer([]byte(customerPayload)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := client.Do(req1)
	if resp1.StatusCode != http.StatusCreated {
		t.Errorf("got %v want %v", resp1.StatusCode, http.StatusCreated)
	}

	// select all customers (will be the newly created one since no other exist)
	var customers []Customer
	req2, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/api/customer", nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("failed to fetch customers: %v", err)
	}
	defer resp2.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	err = json.Unmarshal([]byte(body), &customers)
	if err != nil {
		t.Fatalf("failed to unmarshal response body into costumer struct: %v", err)
	}

	customerID := customers[0].ID

	// restaurant selectitem
	payload := fmt.Sprintf(`{
  		"customerId": %d,
  		"id": 1,
  		"quantity": 2,
  		"restaurantId": 1
	}
	`, customerID)
	req3, _ := http.NewRequest(http.MethodPost, "http://localhost:8083/api/restaurants/menu/select", bytes.NewBuffer([]byte(payload)))
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("got %v want %v", resp3.StatusCode, http.StatusOK)
	}
	defer resp3.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp3.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	responseStr := string(body)
	expectedMessage := `{"message": "Menu item selected successfully"}`

	if responseStr != expectedMessage {
		t.Fatalf("Unexpected response body: got %s, want %s", responseStr, expectedMessage)
	}

	// shoppingcart consume
	req4, _ := http.NewRequest(http.MethodGet, "http://localhost:8084/api/shopping/consume", nil)
	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("got %v want %v", resp4.StatusCode, http.StatusOK)
	}
	defer resp4.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp4.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	responseStr = string(body)
	expectedMessage = `{"message": "Menu item added to cart successfully"}`

	if responseStr != expectedMessage {
		t.Fatalf("Unexpected response body: got %s, want %s", responseStr, expectedMessage)
	}

	// shoppingcart publish
	payload = `{
	  "comment": "No vegetables on the pizza."
	}`
	req5, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8084/api/shopping/publish/%d", customerID), bytes.NewBuffer([]byte(payload)))
	resp5, err := client.Do(req5)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if resp5.StatusCode != http.StatusOK {
		t.Errorf("got %v want %v", resp5.StatusCode, http.StatusOK)
	}
	defer resp5.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp5.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	responseStr = string(body)
	expectedMessage = `{"message": "Shopping Cart published and selected to Order successfully"}`

	if responseStr != expectedMessage {
		t.Fatalf("Unexpected response body: got %s, want %s", responseStr, expectedMessage)
	}

	// order consume
	req6, _ := http.NewRequest(http.MethodGet, "http://localhost:8082/api/order/consume", nil)
	resp6, err := client.Do(req6)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if resp6.StatusCode != http.StatusOK {
		t.Errorf("got %v want %v", resp6.StatusCode, http.StatusOK)
	}
	defer resp6.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp6.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	responseStr = string(body)
	expectedMessage = `{"message": "Order consumption started"}`

	if responseStr != expectedMessage {
		t.Fatalf("Unexpected response body: got %s, want %s", responseStr, expectedMessage)
	}

	// check newly created order
	time.Sleep(5 * time.Second) // maybe not needed
	var orders []generated.Order
	req7, _ := http.NewRequest(http.MethodGet, "http://localhost:8082/api/orders", nil)
	resp7, err := client.Do(req7)
	if err != nil {
		t.Fatalf("failed to fetch orders: %v", err)
	}
	defer resp7.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp7.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Unmarshal the JSON response into the orders slice
	err = json.Unmarshal(body, &orders)
	if err != nil {
		t.Fatalf("failed to unmarshal response body into order struct: %v", err)
	}

	// Access the first order's ID (or any other field)
	if len(orders) == 0 {
		t.Fatalf("no orders found in response")
	}

	order := orders[0]

	// fmt.Println("order", order)
	orderJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal order: %v\n", err)
	} else {
		fmt.Println("Order:", string(orderJSON))
	}
	// assertions

	// // calculate bonus for order
	orderId := orders[0].ID
	// req8, _ := http.NewRequest(http.MethodPost, "", nil)
	//

	// check order again by id
	// // assert that the created order looks right
	// req8.Body

	// cleanup
	// delete customer, shoppingcart and order
	req9, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8081/api/customer/%d", customerID), nil)
	req10, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8084/api/shopping/%d", customerID), nil)
	req11, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8082/api/orders/%d", orderId), nil)
	_, err = client.Do(req9)
	if err != nil {
		t.Fatalf("failed to delete customer: %v", err)
	}
	_, err = client.Do(req10)
	if err != nil {
		t.Fatalf("failed to delete shopping cart: %v", err)
	}
	_, err = client.Do(req11)
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}
}
