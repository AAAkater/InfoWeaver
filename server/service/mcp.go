package service

import (
	"context"
	"server/db"
	"server/models"

	"gorm.io/gorm"
)

var McpServiceApp = new(McpService)

type McpService struct{}

// CreateMcp creates a new MCP server configuration
func (s *McpService) CreateMcp(ctx context.Context, ownerID uint, name string, transport string, command string, args string, url string, headers map[string]any, envVars map[string]any, enabled bool) error {
	dbMcp := &models.Mcp{
		OwnerID:   ownerID,
		Name:      name,
		Transport: transport,
		Command:   command,
		Args:      args,
		URL:       url,
		Headers:   headers,
		EnvVars:   envVars,
		Enabled:   enabled,
	}
	return gorm.G[models.Mcp](db.PgSqlDB).Create(ctx, dbMcp)
}

// UpdateMcp updates an existing MCP server configuration
func (s *McpService) UpdateMcp(ctx context.Context, mcpID uint, ownerID uint, name string, transport string, command string, args string, url string, headers map[string]any, envVars map[string]any, enabled bool) error {
	newMcp := models.Mcp{
		Name:      name,
		Transport: transport,
		Command:   command,
		Args:      args,
		URL:       url,
		Headers:   headers,
		EnvVars:   envVars,
		Enabled:   enabled,
	}
	rows, err := gorm.G[models.Mcp](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", mcpID, ownerID).
		Updates(ctx, newMcp)
	if rows == 0 {
		return ErrNotFound
	}
	return err
}

// GetMcpByID retrieves an MCP server by ID and owner ID
func (s *McpService) GetMcpByID(ctx context.Context, mcpID uint, ownerID uint) (*models.McpInfo, error) {
	var dbMcp models.McpInfo
	result := db.PgSqlDB.Model(&models.Mcp{}).
		Where("id = ? AND owner_id = ?", mcpID, ownerID).
		First(&dbMcp)
	if result.Error != nil {
		return nil, result.Error
	}
	return &dbMcp, nil
}

// CheckMcpExistsByName checks if an MCP server with the given name already exists for the owner
func (s *McpService) CheckMcpExistsByName(ctx context.Context, ownerID uint, name string) (exists bool, err error) {
	cnt, err := gorm.G[models.Mcp](db.PgSqlDB).
		Where("name = ? AND owner_id = ?", name, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}

// CheckMcpOwnership verifies that the MCP server belongs to the specified owner
func (s *McpService) CheckMcpOwnership(ctx context.Context, mcpID uint, ownerID uint) (belongs bool, err error) {
	cnt, err := gorm.G[models.Mcp](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", mcpID, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}

// GetAllMcps retrieves all MCP servers for a specific owner
func (s *McpService) GetAllMcps(ctx context.Context, ownerID uint) (total int64, mcps []models.McpInfo, err error) {
	result := db.PgSqlDB.Model(&models.Mcp{}).
		Where("owner_id = ?", ownerID).
		Find(&mcps)
	return result.RowsAffected, mcps, result.Error
}

// DeleteMcp deletes an MCP server by ID and owner ID
func (s *McpService) DeleteMcp(ctx context.Context, mcpID uint, ownerID uint) error {
	rows, err := gorm.G[models.Mcp](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", mcpID, ownerID).
		Delete(ctx)
	if rows == 0 {
		return ErrNotFound
	}
	return err
}
