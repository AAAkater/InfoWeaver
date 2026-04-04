package service

import (
	"context"
	"fmt"
	"net/url"
	"server/db"
	"server/models"
	"server/utils"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicOption "github.com/anthropics/anthropic-sdk-go/option"
	openaiOption "github.com/openai/openai-go/v3/option"

	"github.com/ollama/ollama/api"
	"github.com/openai/openai-go/v3"
	"google.golang.org/genai"

	"gorm.io/gorm"
)

var ProviderServiceApp = new(ProviderService)

type ProviderService struct{}

func (this *ProviderService) CreateProvider(ctx context.Context, ownerID uint, name string, baseURL string, apiKey string, mode string) error {
	// Encrypt API key using AES instead of bcrypt hash
	encryptedKey, err := utils.EncryptAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	dbProvider := &models.Provider{
		OwnerID: ownerID,
		Name:    name,
		BaseURL: baseURL,
		APIKey:  encryptedKey,
		Mode:    mode,
	}
	return gorm.G[models.Provider](db.PgSqlDB).Create(ctx, dbProvider)
}

func (this *ProviderService) UpdateProvider(ctx context.Context, providerID uint, ownerID uint, name string, baseURL string, apiKey string, mode string) error {
	// Encrypt API key using AES instead of bcrypt hash
	encryptedKey, err := utils.EncryptAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	newProvider := models.Provider{
		Name:    name,
		BaseURL: baseURL,
		APIKey:  encryptedKey,
		Mode:    mode,
	}
	rows, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		Updates(ctx, newProvider)
	if rows == 0 {
		return ErrNotFound
	}
	return err
}
func (this *ProviderService) GetProviderByID(ctx context.Context, providerID uint, ownerID uint) (*models.ProviderInfo, error) {
	var dbProvider models.ProviderInfo
	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		First(&dbProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &dbProvider, nil
}

func (this *ProviderService) GetProviderByName(ctx context.Context, name string, ownerID uint) (*models.ProviderInfo, error) {
	var dbProvider models.ProviderInfo
	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("name = ? AND owner_id = ?", name, ownerID).
		First(&dbProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &dbProvider, nil
}

func (this *ProviderService) CheckProviderExistsByName(ctx context.Context, ownerID uint, name string) (exists bool, err error) {
	cnt, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("name = ? AND owner_id = ?", name, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}

// CheckProviderOwnership verifies that the provider belongs to the specified owner
func (this *ProviderService) CheckProviderOwnership(ctx context.Context, providerID uint, ownerID uint) (belongs bool, err error) {
	cnt, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}

func (this *ProviderService) GetAllProviders(ctx context.Context, ownerID uint) (cows int64, dbProviders []models.ProviderInfo, err error) {

	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("owner_id = ?", ownerID).
		Find(&dbProviders)

	return result.RowsAffected, dbProviders, result.Error
}

func (this *ProviderService) DeleteProvider(ctx context.Context, providerID uint, ownerID uint) error {
	_, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		Delete(ctx)
	return err
}

// GetProviderRawByID retrieves the full provider record (including encrypted API key)
func (this *ProviderService) GetProviderRawByID(ctx context.Context, providerID uint, ownerID uint) (*models.Provider, error) {
	var provider models.Provider
	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		First(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &provider, nil
}

// ListModels fetches available embedding models from the provider's API
func (this *ProviderService) ListModels(ctx context.Context, providerID uint, ownerID uint) (*models.ProviderModelsResp, error) {
	// Get provider with encrypted API key
	provider, err := this.GetProviderRawByID(ctx, providerID, ownerID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Decrypt API key
	apiKey, err := utils.DecryptAPIKey(provider.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// Call external API based on provider mode
	switch provider.Mode {
	case models.PROVIDER_MODE_OPENAI, models.PROVIDER_MODE_OPENAI_RESP:
		return this.listOpenAIModels(ctx, provider.BaseURL, apiKey)
	case models.PROVIDER_MODE_ANTHROPIC:
		return this.listAnthropicModels(ctx, provider.BaseURL, apiKey)
	case models.PROVIDER_MODE_GEMINI:
		return this.listGeminiModels(ctx, provider.BaseURL, apiKey)
	case models.PROVIDER_MODE_OLLAMA:
		return this.listOllamaModels(ctx, provider.BaseURL)
	default:
		return nil, fmt.Errorf("unsupported provider mode: %s", provider.Mode)
	}
}

// listOpenAIModels uses OpenAI SDK to list embedding models
func (this *ProviderService) listOpenAIModels(ctx context.Context, baseURL string, apiKey string) (*models.ProviderModelsResp, error) {
	// Validate that baseURL ends with /v1
	baseURL = strings.TrimSuffix(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		return nil, fmt.Errorf("OpenAI base URL must end with /v1, got: %s", baseURL)
	}

	// Create OpenAI client with custom base URL
	client := openai.NewClient(
		openaiOption.WithAPIKey(apiKey),
		openaiOption.WithBaseURL(baseURL),
	)

	// List all models
	modelsList, err := client.Models.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list OpenAI models: %w", err)
	}

	// Return all models
	allModels := []models.ModelInfo{}
	for _, model := range modelsList.Data {
		allModels = append(allModels, models.ModelInfo{
			ID:      model.ID,
			Object:  string(model.Object),
			OwnedBy: model.OwnedBy,
		})
	}

	return &models.ProviderModelsResp{Models: allModels}, nil
}

// listAnthropicModels uses Anthropic SDK to list available models
func (this *ProviderService) listAnthropicModels(ctx context.Context, baseURL string, apiKey string) (*models.ProviderModelsResp, error) {
	client := anthropic.NewClient(
		anthropicOption.WithAPIKey(apiKey),
		anthropicOption.WithBaseURL(baseURL))

	// List models
	page, err := client.Models.List(ctx, anthropic.ModelListParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Anthropic models: %w", err)
	}

	// Return all models
	allModels := []models.ModelInfo{}
	for _, model := range page.Data {
		allModels = append(allModels, models.ModelInfo{
			ID:      model.ID,
			Object:  string(model.Type),
			OwnedBy: "anthropic",
		})
	}

	return &models.ProviderModelsResp{Models: allModels}, nil
}

// listGeminiModels uses Google Genai SDK to list embedding models
func (this *ProviderService) listGeminiModels(ctx context.Context, baseURL, apiKey string) (*models.ProviderModelsResp, error) {
	// Create Gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		HTTPOptions: genai.HTTPOptions{
			BaseURL: baseURL,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// List models - returns Page[Model] with Items array
	modelsPage, err := client.Models.List(ctx, &genai.ListModelsConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Gemini models: %w", err)
	}

	// Return all models
	allModels := []models.ModelInfo{}
	for _, model := range modelsPage.Items {
		// Extract model name from full path (e.g., "models/gemini-pro")
		name := strings.TrimPrefix(model.Name, "models/")
		allModels = append(allModels, models.ModelInfo{
			ID:      name,
			Object:  "model",
			OwnedBy: "google",
		})
	}

	return &models.ProviderModelsResp{Models: allModels}, nil
}

// listOllamaModels uses Ollama SDK to list embedding models
func (this *ProviderService) listOllamaModels(ctx context.Context, baseURL string) (*models.ProviderModelsResp, error) {
	// Parse base URL
	ollamaURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ollama base URL: %w", err)
	}

	// Create Ollama client
	client := api.NewClient(ollamaURL, nil)

	// List models
	modelsList, err := client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}

	// Return all models
	allModels := []models.ModelInfo{}
	for _, model := range modelsList.Models {
		allModels = append(allModels, models.ModelInfo{
			ID:      model.Name,
			Object:  "model",
			OwnedBy: "ollama",
		})
	}

	return &models.ProviderModelsResp{Models: allModels}, nil
}
