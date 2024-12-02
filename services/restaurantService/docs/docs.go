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
        "/api/categories": {
            "get": {
                "description": "Fetches a list of all unique categories from the restaurant",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Category(Restaurant) CRUD"
                ],
                "summary": "Get all categories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
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
        "/api/filter/{category}": {
            "get": {
                "description": "Fetches all restaurants for a given category",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Category(Restaurant) CRUD"
                ],
                "summary": "Filter restaurants by category",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant Category",
                        "name": "category",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/generated.Restaurant"
                            }
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
        "/api/restaurants": {
            "get": {
                "description": "Fetches a list of all restaurants from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Restaurant CRUD"
                ],
                "summary": "Get all restaurants",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/generated.Restaurant"
                            }
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
        "/api/restaurants/menu/select": {
            "post": {
                "description": "Customer selects a MenuItem or more",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Restaurant Broker"
                ],
                "summary": "Selecting MenuItems",
                "parameters": [
                    {
                        "description": "Menu item selection details",
                        "name": "customer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.SelectItemParams"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Menu item successfully selected",
                        "schema": {
                            "$ref": "#/definitions/handlers.MenuItemSelection"
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
        "/api/restaurants/{id}": {
            "get": {
                "description": "Fetches a restaurant based on the id from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Restaurant CRUD"
                ],
                "summary": "Get restaurant by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/generated.Restaurant"
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
        "/api/restaurants/{restaurantId}/menu-items": {
            "get": {
                "description": "Fetches all menu items associated with a specific restaurant ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MenuItem(Restaurant) CRUD"
                ],
                "summary": "Get menu items by restaurant ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant ID",
                        "name": "restaurantId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/generated.Menuitem"
                            }
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
        "/api/restaurants/{restaurantId}/menu-items/{menuitemId}": {
            "get": {
                "description": "Fetches a menu item based on the restaurant and id from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MenuItem(Restaurant) CRUD"
                ],
                "summary": "Get menu item by restaurant and id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant ID",
                        "name": "restaurantId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Menu Item ID",
                        "name": "menuitemId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/generated.Menuitem"
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
        "generated.Menuitem": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "restaurantid": {
                    "type": "integer"
                }
            }
        },
        "generated.Restaurant": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "zip_code": {
                    "type": "integer"
                }
            }
        },
        "handlers.MenuItemSelection": {
            "type": "object",
            "properties": {
                "customerId": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Cheese Burger"
                },
                "price": {
                    "type": "number",
                    "example": 10
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "restaurantId": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "handlers.SelectItemParams": {
            "type": "object",
            "properties": {
                "customerId": {
                    "type": "integer",
                    "example": 1
                },
                "id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "restaurantId": {
                    "type": "integer",
                    "example": 10
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8083",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Restaurant Service API",
	Description:      "This is the API documentation for the Restaurant Service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
