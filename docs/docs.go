// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "2. Auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Register user request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.registerUserPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "Sign in a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "2. Auth"
                ],
                "summary": "Sign in user",
                "parameters": [
                    {
                        "description": "Sign in user request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.signInPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token",
                        "schema": {
                            "$ref": "#/definitions/main.signInResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/customers": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new customer for store by admin user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "3. Customers"
                ],
                "summary": "Create customer",
                "parameters": [
                    {
                        "description": "Create customer request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.createCustomerPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check if the server is running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "1. Health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.healthCheckResponse"
                        }
                    }
                }
            }
        },
        "/rentals/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a rental by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "4. Rentals"
                ],
                "summary": "Get rental by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Rental ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.rentalResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.createCustomerPayload": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "store_id"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3,
                    "example": "John"
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3,
                    "example": "Doe"
                },
                "store_id": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 1
                }
            }
        },
        "main.healthCheckData": {
            "type": "object",
            "properties": {
                "environment": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "main.healthCheckResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/main.healthCheckData"
                }
            }
        },
        "main.registerUserPayload": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 8,
                    "example": "password123"
                },
                "username": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3,
                    "example": "john.doe"
                }
            }
        },
        "main.rentalResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/store.Rental"
                }
            }
        },
        "main.signInPayload": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 8,
                    "example": "password123"
                }
            }
        },
        "main.signInResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                }
            }
        },
        "store.Rental": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "rental_date": {
                    "type": "string"
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger Examasdasdasdasdasdawdasple API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
