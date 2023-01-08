// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/v1/families": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "family"
                ],
                "summary": "find all families",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.FamiliesResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "family"
                ],
                "summary": "create an family",
                "parameters": [
                    {
                        "description": "Create family",
                        "name": "family",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.FamilyCreateDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/service.FamilyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/families/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "family"
                ],
                "summary": "find family by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "family ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.FamiliesResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "family"
                ],
                "summary": "delete an family",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "family ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "family"
                ],
                "summary": "update an family",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "family ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update family",
                        "name": "family",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.FamilyUpdateDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/persons": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "find all persons",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.PersonResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "create a person",
                "parameters": [
                    {
                        "description": "Create person",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.PersonCreateDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/service.PersonResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/persons/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "find person by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.PersonsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "delete a person",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "update a person",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update person",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.PersonUpdateDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/resources": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "find resource by id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.ResourceResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "create a resource",
                "parameters": [
                    {
                        "description": "Create resource",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.CreateResourceDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/service.ResourceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/resources/{id}": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "update a resource",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update resource",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UpdateResourceDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/resources/{id}/donate": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "donate a resource",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Donate a resource",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.DonateResourceDonateDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/resources/{id}/quantity": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "update a resource quantity",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update resource quantity",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UpdateResourceQuantityDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        },
        "/api/v1/resources/{id}/return": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resource"
                ],
                "summary": "Return a doneted resource",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Return a doneted resource",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.DonateResourceDonateDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.HttpError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "invalid parameter"
                }
            }
        },
        "service.CreateResourceDto": {
            "type": "object",
            "required": [
                "amount",
                "measurement",
                "name",
                "quantity"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0,
                    "example": 5
                },
                "measurement": {
                    "type": "string",
                    "example": "Kg"
                },
                "name": {
                    "type": "string",
                    "example": "Arroz"
                },
                "quantity": {
                    "type": "number",
                    "minimum": 0,
                    "example": 10
                }
            }
        },
        "service.DonateResourceDonateDto": {
            "type": "object",
            "required": [
                "family_id",
                "quantity"
            ],
            "properties": {
                "family_id": {
                    "type": "integer",
                    "example": 1
                },
                "quantity": {
                    "type": "number",
                    "minimum": 0,
                    "example": 10
                }
            }
        },
        "service.FamiliesResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.Family"
                    }
                }
            }
        },
        "service.Family": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string",
                    "example": "São Paulo"
                },
                "complement": {
                    "type": "string",
                    "example": "1A"
                },
                "country": {
                    "type": "string",
                    "example": "BR"
                },
                "created_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "deleted_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Sauro"
                },
                "neighborhood": {
                    "type": "string",
                    "example": "Centro Histórico"
                },
                "number": {
                    "type": "string",
                    "example": "1000"
                },
                "state": {
                    "type": "string",
                    "example": "SP"
                },
                "street": {
                    "type": "string",
                    "example": "R. Vinte e Cinco de Março"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "zipcode": {
                    "type": "string",
                    "example": "01021100"
                }
            }
        },
        "service.FamilyCreateDto": {
            "type": "object",
            "required": [
                "city",
                "complement",
                "country",
                "name",
                "neighborhood",
                "number",
                "state",
                "street",
                "zipcode"
            ],
            "properties": {
                "city": {
                    "type": "string",
                    "example": "São Paulo"
                },
                "complement": {
                    "type": "string",
                    "example": "1A"
                },
                "country": {
                    "type": "string",
                    "example": "BR"
                },
                "name": {
                    "type": "string",
                    "example": "Sauro"
                },
                "neighborhood": {
                    "type": "string",
                    "example": "Centro Histórico"
                },
                "number": {
                    "type": "string",
                    "example": "1000"
                },
                "state": {
                    "type": "string",
                    "example": "SP"
                },
                "street": {
                    "type": "string",
                    "example": "R. Vinte e Cinco de Março"
                },
                "zipcode": {
                    "type": "string",
                    "example": "01021100"
                }
            }
        },
        "service.FamilyResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/service.Family"
                }
            }
        },
        "service.FamilyUpdateDto": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string",
                    "example": "São Paulo"
                },
                "complement": {
                    "type": "string",
                    "example": "1A"
                },
                "country": {
                    "type": "string",
                    "example": "BR"
                },
                "name": {
                    "type": "string",
                    "example": "Sauro"
                },
                "neighborhood": {
                    "type": "string",
                    "example": "Centro Histórico"
                },
                "number": {
                    "type": "string",
                    "example": "1000"
                },
                "state": {
                    "type": "string",
                    "example": "SP"
                },
                "street": {
                    "type": "string",
                    "example": "R. Vinte e Cinco de Março"
                },
                "zipcode": {
                    "type": "string",
                    "example": "01021100"
                }
            }
        },
        "service.Person": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "deleted_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "family_id": {
                    "type": "integer",
                    "example": 1
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Cláudio"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                }
            }
        },
        "service.PersonCreateDto": {
            "type": "object",
            "required": [
                "family_id",
                "name"
            ],
            "properties": {
                "family_id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Cláudio"
                }
            }
        },
        "service.PersonResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/service.Person"
                }
            }
        },
        "service.PersonUpdateDto": {
            "type": "object",
            "properties": {
                "family_id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Cláudio"
                }
            }
        },
        "service.PersonsResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.Person"
                    }
                }
            }
        },
        "service.Resource": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 5
                },
                "created_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                },
                "deleted_at": {
                    "type": "string",
                    "example": ""
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "measurement": {
                    "type": "string",
                    "example": "Kg"
                },
                "name": {
                    "type": "string",
                    "example": "Arroz"
                },
                "quantity": {
                    "type": "number",
                    "example": 10
                },
                "updated_at": {
                    "type": "string",
                    "example": "2000-01-01T12:03:00"
                }
            }
        },
        "service.ResourceResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/service.Resource"
                }
            }
        },
        "service.ResourcesResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.Resource"
                    }
                }
            }
        },
        "service.UpdateResourceDto": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0,
                    "example": 5
                },
                "measurement": {
                    "type": "string",
                    "example": "Kg"
                },
                "name": {
                    "type": "string",
                    "example": "Arroz"
                }
            }
        },
        "service.UpdateResourceQuantityDto": {
            "type": "object",
            "required": [
                "quantity"
            ],
            "properties": {
                "quantity": {
                    "type": "number",
                    "minimum": 0,
                    "example": 10
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
