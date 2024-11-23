// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/shopping": {
            "post": {
                "description": "Add an item to the shopping cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shoppingCart"
                ],
                "summary": "Add an item",
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
        "/api/shopping/{customerId}": {
            "delete": {
                "description": "Clears the cart for the specified customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shoppingCart"
                ],
                "summary": "Clears the cart",
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
                        "description": "cart cleared",
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
                    "shoppingCart"
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
                        "description": "Item updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Item not found",
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
                "description": "Fetches a list of items based on the customerId",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shoppingCart"
                ],
                "summary": "View items for a customer",
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
                        "description": "Cart cleared",
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
                "customer_id": {
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
                "restaurant_id": {
                    "type": "integer"
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
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
