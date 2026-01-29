package models

import (
	"gorm.io/gorm"
)

type (
	// User represents a system user with role-based permissions
	User struct {
		gorm.Model
		Username string `gorm:"unique;not null"`
		Email    string `gorm:"unique;not null"`
		Password string `gorm:"not null"`
		Role     string `gorm:"default:user"` // "user" or "admin"
	}

	// File represents uploaded files stored in MinIO
	File struct {
		gorm.Model
		Name      string `gorm:"not null"`        // Original filename
		MinioPath string `gorm:"not null;unique"` // MinIO object path (bucket/key)
		Size      int64  `gorm:"not null"`        // File size in bytes
		Type      string `gorm:"not null"`        // MIME type
		UserID    uint   `gorm:"not null"`        // Owner user ID
		User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	}

	// Document represents a knowledge document for RAG system
	Document struct {
		gorm.Model
		Content  string         `gorm:"type:text;not null"` // Document content
		Metadata map[string]any `gorm:"type:jsonb"`         // Additional metadata (source, type, etc.)
	}
	// Embedding references vector representations in Milvus
	Embedding struct {
		ID             uint   `gorm:"primaryKey"`
		DocumentID     uint   `gorm:"uniqueIndex;not null"` // Unique embedding per document
		MilvusID       string `gorm:"unique;not null"`      // Reference to Milvus vector ID
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
)
