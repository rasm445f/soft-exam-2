// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/shopping": {
            "post": {
                "description": "Add a MenuItem to the shopping cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart CRUD"
                ],
                "summary": "Add a MenuItem",
                "parameters": [
                    {
                        "description": "item object",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.AddItemParams"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/db.AddItemParams"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/shopping/consume": {
            "get": {
                "description": "Consumes the Shopping Cart's Menu Items for a Customer",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart Broker"
                ],
                "summary": "Consume the chosen Menu Items for a Customer",
                "responses": {
                    "200": {
                        "description": "Shopping Cart's Menu Items Consumed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/shopping/publish/{customerId}": {
            "post": {
                "description": "Selecting the cart for the specified customer with an optional comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart Broker"
                ],
                "summary": "Publish a Customer's shopping cart to RabbitMQ to be consumed by the Order service with an optional Comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Customer ID",
                        "name": "customerId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Customer Comment (optional)",
                        "name": "comment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.PublishShoppingCartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order Selected Successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/shopping/{customerId}": {
            "delete": {
                "description": "Clears the ShoppingCart for a specific customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart CRUD"
                ],
                "summary": "Clears the ShoppingCart",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "customer ID",
                        "name": "customerId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cart Cleared",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/shopping/{customerId}/{itemId}": {
            "patch": {
                "description": "Update the quantity of an existing item in the shopping cart. If the quantity is set to 0, the item will be removed.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart CRUD"
                ],
                "summary": "Update an item in the cart",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "customer ID",
                        "name": "customerId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "itemId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New quantity for the item",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.UpdateQuantityRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ShoppingCart updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/shopping/{id}": {
            "get": {
                "description": "Fetches a list of MenuItems for a specific Customer, to view the ShoppingCart",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ShoppingCart CRUD"
                ],
                "summary": "View MenuItems for a customer's ShoppingCart",
                "parameters": [
                    {
                        "type": "string",
                        "description": "customer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Viewed Cart",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.AddItemParams": {
            "type": "object",
            "properties": {
                "customerId": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "restaurantId": {
                    "type": "integer"
                }
            }
        },
        "handlers.PublishShoppingCartRequest": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string",
                    "example": "No vegetables on the pizza."
                }
            }
        },
        "handlers.UpdateQuantityRequest": {
            "type": "object",
            "properties": {
                "quantity": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8084",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Order Shopping cart API",
	Description:      "This is the API documentation for the Shopping cart Service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
