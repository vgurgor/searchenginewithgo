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
  "paths": {
    "/api/v1/contents/search": {
      "get": {
        "summary": "Search contents",
        "parameters": [
          { "name":"q","in":"query","type":"string","required":false },
          { "name":"type","in":"query","type":"string","enum":["video","text"],"required":false },
          { "name":"sort","in":"query","type":"string","enum":["score_desc","score_asc","date_desc","date_asc"],"required":false },
          { "name":"page","in":"query","type":"integer","required":false },
          { "name":"page_size","in":"query","type":"integer","required":false }
        ],
        "responses": {
          "200": { "description":"OK" }
        }
      }
    },
    "/api/v1/contents/{id}": {
      "get": {
        "summary":"Get content by ID",
        "parameters":[
          {"name":"id","in":"path","type":"integer","required":true}
        ],
        "responses": {
          "200": { "description":"OK" },
          "404": { "description":"Not Found" }
        }
      }
    },
    "/api/v1/contents/stats": {
      "get": {
        "summary":"Get contents statistics",
        "responses": {
          "200": { "description":"OK" }
        }
      }
    }
  },
  "definitions": {
    "Content": {
      "type": "object",
      "properties": {
        "id": { "type": "integer", "format": "int64" },
        "providerId": { "type": "string" },
        "providerContentId": { "type": "string" },
        "title": { "type": "string" },
        "contentType": { "type": "string", "enum": ["video","text"] },
        "description": { "type": "string" },
        "url": { "type": "string" },
        "thumbnailUrl": { "type": "string" },
        "publishedAt": { "type": "string", "format": "date-time" },
        "createdAt": { "type": "string", "format": "date-time" },
        "updatedAt": { "type": "string", "format": "date-time" }
      }
    },
    "ContentMetrics": {
      "type": "object",
      "properties": {
        "id": { "type": "integer", "format": "int64" },
        "contentId": { "type": "integer", "format": "int64" },
        "views": { "type": "integer", "format": "int64" },
        "likes": { "type": "integer", "format": "int64" },
        "readingTime": { "type": "integer" },
        "reactions": { "type": "integer" },
        "finalScore": { "type": "number", "format": "double" },
        "recalculatedAt": { "type": "string", "format": "date-time" },
        "createdAt": { "type": "string", "format": "date-time" },
        "updatedAt": { "type": "string", "format": "date-time" }
      }
    },
    "SearchFilters": {
      "type": "object",
      "properties": {
        "keyword": { "type": "string" },
        "contentType": { "type": "string", "enum": ["video","text"] },
        "sortBy": { "type": "string", "enum": ["popularity","relevance"] },
        "page": { "type": "integer", "default": 1 },
        "pageSize": { "type": "integer", "default": 20, "maximum": 100 }
      }
    },
    "ContentResponse": {
      "type": "object",
      "properties": {
        "id": { "type": "integer", "format": "int64" },
        "title": { "type": "string" },
        "contentType": { "type": "string", "enum": ["video","text"] },
        "description": { "type": "string" },
        "url": { "type": "string" },
        "thumbnailUrl": { "type": "string" },
        "score": { "type": "number", "format": "double" },
        "publishedAt": { "type": "string", "format": "date-time" }
      }
    },
    "PaginatedResponse_ContentResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": { "$ref": "#/definitions/ContentResponse" }
        },
        "page": { "type": "integer" },
        "pageSize": { "type": "integer" },
        "totalCount": { "type": "integer", "format": "int64" },
        "totalPages": { "type": "integer" }
      }
    }
  }
}`

type s struct{}

func (s *s) ReadDoc() string { return doc }

func init() {
	swag.Register(swag.Name, &s{})
}


