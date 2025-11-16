package docs

import "github.com/swaggo/swag"

// Minimal embedded Swagger doc so that gin-swagger serves UI at /swagger
var doc = `{
  "swagger": "2.0",
  "info": {
    "description": "Content Search Engine - Public & Admin API",
    "title": "Content Search Engine API",
    "version": "1.0.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "schemes": ["http"],
  "produces": ["application/json"],
  "consumes": ["application/json"],
  "securityDefinitions": {
    "ApiKeyAuth": { "type":"apiKey", "name":"X-API-Key", "in":"header" }
  },
  "paths": {
    "/api/v1/contents/search": {
      "get": {
        "summary": "Search contents",
        "description": "Search and filter contents by keyword, type and sort order with pagination",
        "tags": ["Contents"],
        "parameters": [
          { "name":"q", "in":"query", "type":"string", "required": false, "description": "Search keyword (title, description)", "maxLength": 100, "example":"travel" },
          { "name":"type", "in":"query", "type":"string", "enum":["video","text"], "required": false, "description": "Content type filter", "example":"video" },
          { "name":"sort", "in":"query", "type":"string", "enum":["score_desc","score_asc","date_desc","date_asc"], "required": false, "description": "Sort order", "example":"score_desc" },
          { "name":"page", "in":"query", "type":"integer", "required": false, "description": "Page number (default: 1)", "minimum": 1, "example":1 },
          { "name":"page_size", "in":"query", "type":"integer", "required": false, "description": "Items per page (default: 20, max: 100)", "minimum": 1, "maximum": 100, "example":20 }
        ],
        "responses": {
          "200": { "description":"OK", "schema": { "$ref":"#/definitions/SearchResponse" },
            "examples": { "application/json": {
              "success": true,
              "data": [
                { "id": 1, "title":"Amazing Video Title", "content_type":"video", "description":"This is a video", "url":"https://example.com/v.mp4", "thumbnail_url":"https://example.com/t.jpg", "score":245.67, "published_at":"2024-11-01T10:00:00Z", "provider":"provider1" }
              ],
              "pagination": { "page":1, "page_size":20, "total_items":150, "total_pages":8 }
            } }
          },
          "400": { "description":"Bad Request", "schema": { "$ref":"#/definitions/SearchResponse" },
            "examples": { "application/json": { "success": false, "error": { "code":"INVALID_PARAMETER", "message":"page_size must be between 1 and 100" } } }
          },
          "500": { "description":"Internal Error" }
        }
      }
    },
    "/api/v1/contents/{id}": {
      "get": {
        "summary": "Get content by ID",
        "tags": ["Contents"],
        "parameters": [
          { "name":"id", "in":"path", "type":"integer", "format":"int64", "required": true }
        ],
        "responses": {
          "200": { "description":"OK", "schema": { "$ref":"#/definitions/APIContentResponse" },
            "examples": { "application/json": {
              "success": true,
              "data": { "id":1, "title":"Amazing Video Title", "content_type":"video", "description":"Full description", "url":"https://example.com/v.mp4", "thumbnail_url":"https://example.com/t.jpg", "score":245.67, "published_at":"2024-11-01T10:00:00Z", "provider":"provider1",
                "metrics": { "views":150000, "likes":5000, "reading_time": null, "reactions": null, "recalculated_at":"2024-11-16T10:30:00Z" } }
            } }
          },
          "404": { "description":"Not Found", "schema": { "$ref":"#/definitions/APIContentResponse" },
            "examples": { "application/json": { "success": false, "error": { "code":"CONTENT_NOT_FOUND", "message":"Content not found" } } }
          }
        }
      }
    },
    "/api/v1/contents/stats": {
      "get": {
        "summary": "Get global statistics",
        "tags": ["Contents"],
        "responses": {
          "200": { "description":"OK", "schema": { "$ref":"#/definitions/StatsResponse" },
            "examples": { "application/json": {
              "success": true,
              "data": {
                "total_contents": 1250, "total_videos": 780, "total_texts": 470, "average_score": 138.92,
                "providers": [ { "provider_id": "provider1", "content_count":650 }, { "provider_id": "provider2", "content_count":600 } ]
              }
            } }
          }
        }
      }
    },

    "/api/v1/admin/sync": {
      "post": {
        "summary": "Trigger manual sync",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "parameters": [
          { "in":"body", "name":"body", "schema": { "$ref":"#/definitions/SyncRequest" } }
        ],
        "responses": {
          "200": { "description":"OK",
            "examples": { "application/json": {
              "success": true,
              "data": { "results": [ { "provider_id":"provider1", "total_fetched":45, "new_contents":5, "updated_contents":10, "skipped_contents":30, "failed_contents":0, "duration_ms":2340, "synced_at":"2024-11-16T10:30:15Z" } ] }
            } }
          },
          "202": { "description":"Accepted",
            "examples": { "application/json": { "success": true, "message":"Sync job started", "job_id":"sync-uuid-1234", "data":{ "provider_id":"provider1", "started_at":"2024-11-16T10:30:00Z" } } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/sync/history": {
      "get": {
        "summary": "List sync history",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "parameters": [
          { "name":"provider_id", "in":"query", "type":"string", "required": false },
          { "name":"status", "in":"query", "type":"string", "enum":["success","partial","failed","in_progress","skipped"], "required": false },
          { "name":"limit", "in":"query", "type":"integer", "required": false },
          { "name":"offset", "in":"query", "type":"integer", "required": false }
        ],
        "responses": {
          "200": { "description":"OK", "schema": { "$ref":"#/definitions/SyncHistoryListResponse" },
            "examples": { "application/json": {
              "success": true,
              "data": [ { "id":"uuid-1234","provider_id":"provider1","sync_status":"success","total_fetched":45,"new_contents":5,"updated_contents":10,"skipped_contents":30,"failed_contents":0,"error_message":null,"started_at":"2024-11-16T08:00:00Z","completed_at":"2024-11-16T08:02:20Z","duration_ms":140000 } ],
              "pagination": { "limit":50, "offset":0, "total":150 }
            } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/scores/recalculate": {
      "post": {
        "summary": "Trigger score recalculation",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "parameters": [
          { "in":"body", "name":"body", "schema": { "$ref":"#/definitions/ScoreRecalcRequest" } }
        ],
        "responses": {
          "202": { "description":"Accepted", "schema": { "$ref":"#/definitions/JobResponse" },
            "examples": { "application/json": { "success": true, "message": "Score recalculation job started", "job_id": "recalc-uuid-9999", "data": { "estimated_items": 1250, "started_at": "2024-11-16T10:35:00Z" } } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/providers": {
      "get": {
        "summary": "Providers info & stats",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "responses": {
          "200": { "description":"OK", "schema": { "type":"array", "items": { "$ref":"#/definitions/ProviderInfo" } },
            "examples": { "application/json": {
              "success": true,
              "data": [
                { "provider_id":"provider1","provider_type":"json","base_url":"http://localhost:8080/mock/provider1","rate_limit":100,"status":"active","last_sync":"2024-11-16T08:00:00Z","last_sync_status":"success","content_count":650,"average_score":145.67 }
              ]
            } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/providers/health-check": {
      "post": {
        "summary": "Check providers health",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "responses": {
          "200": { "description":"OK", "schema": { "type":"array", "items": { "$ref":"#/definitions/ProviderHealth" } },
            "examples": { "application/json": {
              "success": true,
              "data": [ { "provider_id":"provider1","is_healthy":true,"response_time_ms":245,"status_code":200,"checked_at":"2024-11-16T10:40:00Z","error":null } ]
            } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/contents/{id}": {
      "delete": {
        "summary": "Soft delete a content",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "parameters": [
          { "name":"id", "in":"path", "type":"integer", "format":"int64", "required": true }
        ],
        "responses": {
          "200": { "description":"OK",
            "examples": { "application/json": { "success": true, "message": "Content deleted successfully", "data": { "id": "uuid-1234", "deleted_at": "2024-11-16T10:45:00Z" } } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/metrics/dashboard": {
      "get": {
        "summary": "Dashboard metrics",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "responses": {
          "200": { "description":"OK",
            "examples": { "application/json": {
              "success": true,
              "data": {
                "overview": { "total_contents": 1250, "total_videos": 780, "total_texts": 470, "average_score": 138.92, "providers_count": 2 },
                "sync_stats": { "last_sync": "2024-11-16T08:00:00Z", "total_syncs": 4 },
                "content_distribution": [ { "provider_id":"provider1","content_count":650,"percentage":52.0 }, { "provider_id":"provider2","content_count":600,"percentage":48.0 } ]
              }
            } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    },
    "/api/v1/admin/jobs/{jobId}": {
      "get": {
        "summary": "Get job status",
        "tags": ["Admin"],
        "security": [ { "ApiKeyAuth": [] } ],
        "parameters": [
          { "name":"jobId", "in":"path", "type":"string", "required": true }
        ],
        "responses": {
          "200": { "description":"OK", "schema": { "$ref":"#/definitions/JobResponse" },
            "examples": { "application/json": {
              "success": true,
              "job_id": "sync-uuid-1234",
              "data": { "job_id":"sync-uuid-1234","type":"sync","status":"running","progress":45,"started_at":"2024-11-16T10:30:00Z","ended_at":null,"error":null }
            } }
          },
          "401": { "description":"Unauthorized" }
        }
      }
    }
  },
  "definitions": {
    "ErrorDTO": {
      "type": "object",
      "properties": {
        "code": { "type":"string", "description":"Machine-readable error code" },
        "message": { "type":"string", "description":"Human-readable error message" },
        "details": { "type":"object", "additionalProperties": { "type":"string" }, "description":"Optional field-level details" }
      }
    },
    "PaginationDTO": {
      "type": "object",
      "properties": {
        "page": { "type":"integer", "description":"Current page number (1-based)", "minimum": 1 },
        "page_size": { "type":"integer", "description":"Items per page", "minimum": 1, "maximum": 100 },
        "total_items": { "type":"integer", "format":"int64", "description":"Total items matching the query", "minimum": 0 },
        "total_pages": { "type":"integer", "description":"Total pages based on page_size", "minimum": 0 }
      }
    },
    "ContentSummaryDTO": {
      "type": "object",
      "properties": {
        "id": { "type":"integer", "format":"int64", "description":"Internal content ID" },
        "title": { "type":"string", "description":"Content title", "example":"Amazing Video Title" },
        "content_type": { "type":"string", "enum":["video","text"], "description":"Content type", "example":"video" },
        "description": { "type":"string", "description":"Short description (may be truncated)", "maxLength": 200, "example":"This is a video description..." },
        "url": { "type":"string", "description":"Canonical content URL", "example":"https://example.com/video.mp4" },
        "thumbnail_url": { "type":"string", "description":"Thumbnail image URL", "example":"https://example.com/thumb.jpg" },
        "score": { "type":"number", "format":"double", "description":"Computed ranking score (2 decimal precision)", "minimum": 0, "example": 245.67 },
        "published_at": { "type":"string", "format":"date-time", "description":"Original publish timestamp (UTC)", "example":"2024-11-01T10:00:00Z" },
        "provider": { "type":"string", "description":"Provider identifier", "example":"provider1" }
      }
    },
    "MetricsDTO": {
      "type": "object",
      "properties": {
        "views": { "type":"integer", "format":"int64", "description":"View count (video only)", "minimum": 0, "example":150000 },
        "likes": { "type":"integer", "format":"int64", "description":"Like count (video only)", "minimum": 0, "example":5000 },
        "reading_time": { "type":"integer", "description":"Estimated reading time in minutes (text only)", "minimum": 0, "example": 8 },
        "reactions": { "type":"integer", "description":"Reaction count (text only)", "minimum": 0, "example": 250 },
        "recalculated_at": { "type":"string", "format":"date-time", "description":"Last score recalculation time (UTC)", "example":"2024-11-16T10:30:00Z" }
      }
    },
    "ContentDetailDTO": {
      "type": "object",
      "properties": {
        "id": { "type":"integer", "format":"int64", "description":"Internal content ID" },
        "title": { "type":"string", "description":"Content title" },
        "content_type": { "type":"string", "enum":["video","text"], "description":"Content type" },
        "description": { "type":"string", "description":"Full description (if available)" },
        "url": { "type":"string", "description":"Canonical content URL", "example":"https://example.com/video.mp4" },
        "thumbnail_url": { "type":"string", "description":"Thumbnail URL", "example":"https://example.com/thumb.jpg" },
        "score": { "type":"number", "format":"double", "description":"Computed ranking score", "example": 245.67 },
        "published_at": { "type":"string", "format":"date-time", "description":"Publish timestamp (UTC)", "example":"2024-11-01T10:00:00Z" },
        "provider": { "type":"string", "description":"Provider identifier", "example":"provider1" },
        "metrics": { "$ref":"#/definitions/MetricsDTO", "description":"Engagement metrics" }
      }
    },
    "SearchResponse": {
      "type": "object",
      "properties": {
        "success": { "type":"boolean", "description":"Operation success flag" },
        "data": { "type":"array", "items": { "$ref":"#/definitions/ContentSummaryDTO" }, "description":"Result items" },
        "pagination": { "$ref":"#/definitions/PaginationDTO", "description":"Pagination metadata" },
        "error": { "$ref":"#/definitions/ErrorDTO", "description":"Error details when success=false" }
      }
    },
    "APIContentResponse": {
      "type": "object",
      "properties": {
        "success": { "type":"boolean", "description":"Operation success flag" },
        "data": { "$ref":"#/definitions/ContentDetailDTO", "description":"Content details" },
        "error": { "$ref":"#/definitions/ErrorDTO", "description":"Error details when success=false" }
      }
    },
    "StatsProviderDTO": {
      "type": "object",
      "properties": {
        "provider_id": { "type":"string", "description":"Provider identifier" },
        "content_count": { "type":"integer", "format":"int64", "description":"Total contents for provider" },
        "last_sync": { "type":"string", "format":"date-time", "description":"Last successful sync time (if any)" }
      }
    },
    "StatsDTO": {
      "type": "object",
      "properties": {
        "total_contents": { "type":"integer", "format":"int64", "description":"Total content count", "minimum": 0, "example": 1250 },
        "total_videos": { "type":"integer", "format":"int64", "description":"Total videos count", "minimum": 0, "example": 780 },
        "total_texts": { "type":"integer", "format":"int64", "description":"Total texts count", "minimum": 0, "example": 470 },
        "average_score": { "type":"number", "format":"double", "description":"Average final score", "minimum": 0, "example": 138.92 },
        "last_sync": { "type":"string", "format":"date-time", "description":"Latest sync time across providers" },
        "providers": { "type":"array", "items": { "$ref":"#/definitions/StatsProviderDTO" }, "description":"Per provider breakdown" }
      }
    },
    "StatsResponse": {
      "type": "object",
      "properties": {
        "success": { "type":"boolean", "description":"Operation success flag" },
        "data": { "$ref":"#/definitions/StatsDTO", "description":"Stats payload" }
      }
    },
    "SyncRequest": {
      "type": "object",
      "properties": {
        "provider_id": { "type":"string", "description":"Specific provider to sync (optional)", "maxLength": 50, "pattern": "^[a-zA-Z0-9_-]+$", "example":"provider1" },
        "force": { "type":"boolean", "description":"Ignore caches if implemented (optional)" },
        "async": { "type":"boolean", "description":"Run as async job (default true if enabled)" }
      }
    },
    "SyncHistory": {
      "type": "object",
      "properties": {
        "id": { "type":"integer", "format":"int64", "description":"History record ID" },
        "provider_id": { "type":"string", "description":"Provider identifier" },
        "sync_status": { "type":"string", "enum":["success","partial","failed","in_progress","skipped"], "description":"Final status" },
        "total_fetched": { "type":"integer", "description":"Fetched items count" },
        "new_contents": { "type":"integer", "description":"Newly created items" },
        "updated_contents": { "type":"integer", "description":"Updated items" },
        "skipped_contents": { "type":"integer", "description":"Skipped items" },
        "failed_contents": { "type":"integer", "description":"Failed items" },
        "error_message": { "type":"string", "description":"Optional error" },
        "started_at": { "type":"string", "format":"date-time", "description":"Start time (UTC)" },
        "completed_at": { "type":"string", "format":"date-time", "description":"Completion time (UTC)" },
        "duration_ms": { "type":"integer", "description":"Duration in milliseconds" }
      }
    },
    "SyncHistoryListResponse": {
      "type": "object",
      "properties": {
        "success": { "type":"boolean", "description":"Operation success flag" },
        "data": { "type":"array", "items": { "$ref":"#/definitions/SyncHistory" }, "description":"History items" },
        "pagination": {
          "type":"object",
          "properties": {
            "limit": { "type":"integer", "description":"Page size", "minimum": 1, "maximum": 200 },
            "offset": { "type":"integer", "description":"Offset for listing", "minimum": 0 },
            "total": { "type":"integer", "description":"Total history items", "minimum": 0 }
          }
        }
      }
    },
    "ScoreRecalcRequest": {
      "type": "object",
      "properties": {
        "content_id": { "type":"integer", "format":"int64", "description":"Specific content ID", "minimum": 1, "example": 123 },
        "content_type": { "type":"string", "enum":["video","text"], "description":"Recalculate by type", "example":"video" },
        "recalculate_all": { "type":"boolean", "description":"Recalculate for all contents", "example": false }
      }
    },
    "JobInfo": {
      "type": "object",
      "properties": {
        "job_id": { "type":"string", "description":"Job identifier" },
        "type": { "type":"string", "description":"Job type (sync, recalculate)" },
        "status": { "type":"string", "enum":["pending","running","completed","failed"], "description":"Current job status" },
        "progress": { "type":"integer", "description":"Progress percentage" },
        "started_at": { "type":"string", "format":"date-time", "description":"Start time (UTC)" },
        "ended_at": { "type":"string", "format":"date-time", "description":"End time (UTC)" },
        "error": { "type":"string", "description":"Error if failed" }
      }
    },
    "JobResponse": {
      "type": "object",
      "properties": {
        "success": { "type":"boolean", "description":"Operation success flag" },
        "job_id": { "type":"string", "description":"Job identifier" },
        "data": { "$ref":"#/definitions/JobInfo", "description":"Job payload" }
      }
    },
    "ProviderInfo": {
      "type": "object",
      "properties": {
        "provider_id": { "type":"string", "description":"Provider identifier", "example":"provider1" },
        "content_count": { "type":"integer", "format":"int64", "description":"Total contents", "minimum": 0, "example": 650 },
        "average_score": { "type":"number", "format":"double", "description":"Average final score", "minimum": 0, "example": 145.67 },
        "last_sync": { "type":"string", "format":"date-time", "description":"Last sync time (optional)" },
        "last_sync_status": { "type":"string", "description":"Last sync status (optional)" }
      }
    },
    "ProviderHealth": {
      "type": "object",
      "properties": {
        "provider_id": { "type":"string", "description":"Provider identifier", "example":"provider1" },
        "is_healthy": { "type":"boolean", "description":"Provider health boolean", "example": true },
        "response_time_ms": { "type":"integer", "description":"Response time in ms", "minimum": 0, "example": 245 },
        "status_code": { "type":"integer", "description":"HTTP status code", "minimum": 100, "maximum": 599, "example": 200 },
        "checked_at": { "type":"string", "format":"date-time", "description":"Check time (UTC)", "example":"2024-11-16T10:40:00Z" },
        "error": { "type":"string", "description":"Error details if any" }
      }
    }
  }
}`

type s struct{}

func (s *s) ReadDoc() string { return doc }

func init() {
	swag.Register(swag.Name, &s{})
}
