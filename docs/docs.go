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
        "/expression/add": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "expression"
                ],
                "summary": "Add an expression",
                "parameters": [
                    {
                        "description": "expression to calculate",
                        "name": "expression",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ExpressionAdding"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "id of just created expression",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "incorrect body",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "unprocessed error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/expression/all": {
            "get": {
                "tags": [
                    "expression"
                ],
                "summary": "Get all expressions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Expression"
                            }
                        }
                    },
                    "500": {
                        "description": "unprocessed error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/worker/want_to_calculate": {
            "get": {
                "tags": [
                    "worker"
                ],
                "summary": "One available expression for worker",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Expression"
                        }
                    },
                    "404": {
                        "description": "no rows now",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "unprocessed error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Error": {
            "type": "object",
            "properties": {
                "errorMessage": {
                    "type": "string"
                }
            }
        },
        "models.Expression": {
            "type": "object",
            "properties": {
                "answer": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "incomingDate": {
                    "type": "integer"
                },
                "progress": {
                    "type": "string"
                },
                "vanilla": {
                    "type": "string"
                }
            }
        },
        "models.ExpressionAdding": {
            "type": "object",
            "required": [
                "expression"
            ],
            "properties": {
                "expression": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
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
