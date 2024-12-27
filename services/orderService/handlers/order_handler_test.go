package handlers

//TODO: implement system test
// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"
// 	"testing"
// )
//
// func TestSystem(t *testing.T) {
// 	// make sure all the services are running first
// 	client := &http.Client{}
//
// 	// create a customer
// 	customerPayload := `{
// 		"email": "test@test.dk",
// 		"name": "string",
// 		"password": "Test123!!!!",
// 		"phonenumber": "12345678",
// 		"street_address": "ligma",
// 		"zip_code": 2800
// 	}`
//
// 	// Create the request
// 	req1, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/customer", bytes.NewBuffer([]byte(customerPayload)))
// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}
// 	req1.Header.Set("Content-Type", "application/json")
// 	resp1, err := client.Do(req1)
// 	fmt.Println(resp1.StatusCode)
//
// 	// select all customers (will be the newly created one since no other exist)
// 	req2, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/api/customer", nil)
// 	resp, err = client.Do(req2)
// 	if err != nil {
// 		t.Fatalf("failed to fetch customers: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	// Read the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Fatalf("failed to read response body: %v", err)
// 	}
//
// 	// Find the ID in the response body
// 	bodyStr := string(body)
// 	idStart := strings.Index(bodyStr, `"id":"`) + len(`"id":"`)
// 	idEnd := strings.Index(bodyStr[idStart:], `"`) + idStart
// 	customerID := bodyStr[idStart:idEnd]
//
// 	t.Logf("First customer ID: %s", customerID)
//
// 	// // restaurant selectitem
// 	// req3, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // shoppingcart consume
// 	// req4, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // shoppingcart publish
// 	// req5, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // order consume
// 	// req6, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // order calc bonus
// 	// req7, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // get order
// 	// req8, _ := http.NewRequest(http.MethodPost, "", nil)
// 	//
// 	// // assert that the created order looks right
// 	// req8.Body
// }
