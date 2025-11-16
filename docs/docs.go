package docs

import "github.com/swaggo/swag"

// Minimal embedded Swagger doc so that gin-swagger serves UI at /swagger
var doc = `{
  "swagger": "2.0",
  "info": {
    "description": "Search Engine API - Swagger",
    "title": "Search Engine API",
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "paths": {}
}`

type s struct{}

func (s *s) ReadDoc() string { return doc }

func init() {
	swag.Register(swag.Name, &s{})
}


