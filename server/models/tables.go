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

	// Document represents a knowledge document for RAG system
	//
	// file chunks are processed into documents
	Document struct {
		gorm.Model
		Content  string         `gorm:"type:text;not null"` // Document content
		Metadata map[string]any `gorm:"type:jsonb"`         // Additional metadata (source, type, etc.)
	}
	// Embedding references vector representations in Milvus
	Embedding struct {
		ID             uint   `gorm:"primaryKey"`
		DocumentID     uint   `gorm:"uniqueIndex;not null"` // Unique embedding per document
		VectorID       string `gorm:"unique;not null"`      // Reference to Milvus vector ID
		CollectionName string `gorm:"not null"`             // Milvus collection name
	}
	// Memory stores user interaction history and retrieval results
	Memory struct {
		gorm.Model
		Question string `gorm:"type:text;not null"`
		Answer   string `gorm:"type:text"` // Generated answer
		// Many-to-Many relationship with RetrievedDocuments
		RetrievedDocuments []Document `gorm:"many2many:memory_documents;"`
	}

	// Dataset represents a collection of files owned by a user
	Dataset struct {
		gorm.Model
		Name        string `gorm:"not null;unique"`
		Icon        string // Icon is an emoji (e.g., 🚀, ❤️).
		Description string
		OwnerID     uint `gorm:"not null"`
		User        User `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	}
	// Provider represents an AI model provider (OpenAI, Gemini, Anthropic, Ollama, etc.)
	Provider struct {
		gorm.Model
		Name    string `gorm:"not null;unique"` // Provider Name: "openai","deepseek","qwen",  "gemini", "anthropic", "ollama"
		Mode    string `gorm:"not null"`        // Provider mode: "openai","openai response", "gemini", "anthropic", "ollama"
		BaseURL string `gorm:"not null"`        // Base URL for API requests
		APIKey  string `gorm:"not null"`        // API key for authentication
		OwnerID uint   `gorm:"not null"`
		User    User   `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	}
)
