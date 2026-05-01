package models

import (
	"gorm.io/gorm"
)

type (
	// User represents a system user with role-based permissions
	User struct {
		gorm.Model
		Username string `gorm:"not null"`
		Email    string `gorm:"unique;not null"`
		Password string `gorm:"not null"`
		Role     string `gorm:"default:user"` // "user" or "admin"
	}

	// File represents uploaded files stored in MinIO
	File struct {
		gorm.Model
		Name      string  `gorm:"not null"`        // Original filename
		MinioPath string  `gorm:"not null;unique"` // MinIO object path (bucket/key)
		Size      int64   `gorm:"not null"`        // File size in bytes
		Type      string  `gorm:"not null"`        // MIME type
		DatasetID uint    `gorm:"not null"`        // Associated dataset ID
		UserID    uint    `gorm:"not null"`        // Owner user ID
		User      User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
		Dataset   Dataset `gorm:"foreignKey:DatasetID;constraint:OnDelete:CASCADE"`
	}

	// Chunk represents a knowledge document for RAG system
	Chunk struct {
		gorm.Model
		Content       string         `gorm:"type:text;not null"`                        // Document content
		ChunkMetadata map[string]any `gorm:"type:jsonb"`                                // Additional metadata (source, type, etc.)
		Status        string         `gorm:"type:varchar(20);not null;default:pending"` // pending | embedding | completed | failed
		VectorID      string         `gorm:"type:varchar(255);uniqueIndex;index"`       // Reference to Milvus vector ID (nullable until embedding completes)
		FileID        uint           `gorm:"not null"`                                  // Source file ID
		File          File           `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE"`
	}
	// Memory stores a single chat message within a session
	Memory struct {
		gorm.Model
		SessionID          uint        `gorm:"not null;index"`              // Associated chat session ID
		Content            string      `gorm:"type:text;not null"`          // Message content
		Role               string      `gorm:"not null;default:'user'"`     // Message role: 'user', 'assistant', 'system'
		RetrievedDocuments []Chunk     `gorm:"many2many:memory_documents;"` // Retrieved chunks for RAG
		ChatSession        ChatSession `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	}

	// Dataset represents a collection of files owned by a user
	Dataset struct {
		gorm.Model
		Name           string `gorm:"not null"`
		Icon           string // Icon is an emoji (e.g., 🚀, ❤️).
		Description    string
		SearchType     string   `gorm:"not null;default:'dense'"` // "sparse", "dense", "hybrid"
		EmbeddingModel string   `gorm:"not null"`                 // Embedding model name (required)
		ProviderID     uint     `gorm:"not null"`                 // Associated provider ID for API access
		OwnerID        uint     `gorm:"not null"`
		User           User     `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
		Provider       Provider `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
	}
	// Provider represents an AI model provider (OpenAI, Gemini, Anthropic, Ollama, etc.)
	Provider struct {
		gorm.Model
		Name            string         `gorm:"not null;"`  // Provider Name
		Mode            string         `gorm:"not null"`   // Provider mode: "openai","openai response", "gemini", "anthropic", "ollama"
		BaseURL         string         `gorm:"not null"`   // Base URL for API requests
		APIKey          string         `gorm:"not null"`   // API key for authentication
		AvailableModels ModelEnableMap `gorm:"type:jsonb"` // Available models with enable status (model_id -> enabled)
		OwnerID         uint           `gorm:"not null"`
		User            User           `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	}

	// ChatSession represents a conversation session with an AI model
	ChatSession struct {
		gorm.Model
		Title   string `gorm:"not null"` // Session title
		OwnerID uint   `gorm:"not null"` // Owner user ID
		User    User   `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	}

	// Mcp represents a Model Context Protocol server configuration
	Mcp struct {
		gorm.Model
		Name      string         `gorm:"not null"`                 // Display name for this MCP server
		Transport string         `gorm:"not null;default:'stdio'"` // Transport type: "stdio", "sse", "streamable_http"
		Command   string         // Command to run (for stdio transport)
		Args      string         `gorm:"type:text"` // JSON-encoded array of command arguments
		URL       string         // Server URL (for sse/streamable_http transport)
		Headers   map[string]any `gorm:"type:jsonb"`   // Custom HTTP headers
		EnvVars   map[string]any `gorm:"type:jsonb"`   // Environment variables for the command
		Enabled   bool           `gorm:"default:true"` // Whether this MCP server is active
		OwnerID   uint           `gorm:"not null"`     // Owner user ID
		User      User           `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	}
)
