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
            "name": "Varun Rajput",
            "url": "https://github.com/rvarun11",
            "email": "rajpuv2@mcmaster.ca"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/player": {
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "Update Player",
                "operationId": "update-challenge",
                "parameters": [
                    {
                        "description": "Player data",
                        "name": "player",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.playerDTO"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/api/v1/player/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "Get Player by ID",
                "operationId": "get-player",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player ID (UUID)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.playerDTO"
                        }
                    }
                }
            }
        },
        "/api/v1/players": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "List Players",
                "operationId": "list-players",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/http.playerDTO"
                            }
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "Create a Player",
                "operationId": "create-player",
                "parameters": [
                    {
                        "description": "Player data",
                        "name": "player",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.playerDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.playerDTO"
                        }
                    }
                }
            }
        },
        "/api/v1/players/{id}": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "Delete a Player by ID",
                "operationId": "delete-player",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player ID (UUID) to be deleted",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        }
    },
    "definitions": {
        "http.playerDTO": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "CreatedAt is the time when the player registered",
                    "type": "string"
                },
                "height": {
                    "description": "Height of the player",
                    "type": "number"
                },
                "id": {
                    "description": "ID of the player",
                    "type": "string"
                },
                "preference": {
                    "description": "Preference of the player",
                    "type": "string"
                },
                "updated_at": {
                    "description": "UpdatedAt is the time when the player last updated the profile",
                    "type": "string"
                },
                "user": {
                    "description": "User is the root entity of player",
                    "allOf": [
                        {
                            "$ref": "#/definitions/http.userDTO"
                        }
                    ]
                },
                "weight": {
                    "description": "Weight of the player",
                    "type": "number"
                },
                "zone_id": {
                    "description": "GeographicalZone is a group of trails in a region",
                    "type": "string"
                }
            }
        },
        "http.userDTO": {
            "type": "object",
            "properties": {
                "dob": {
                    "description": "DoB",
                    "type": "string"
                },
                "email": {
                    "description": "Email",
                    "type": "string"
                },
                "id": {
                    "description": "ID is the identifier of the Entity, the ID is shared for all sub domains",
                    "type": "string"
                },
                "name": {
                    "description": "Name of the user",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "User Manager API",
	Description:      "This provides a description of API endpoints for the Challenge Manager",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}