package utils

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var validatorInstance = validator.New(validator.WithRequiredStructEnabled())

func init() {
	// Register custom validator
	_ = validatorInstance.RegisterValidation("emoji", validateEmoji)
}

func BindAndValidate[T any](ctx *echo.Context) (*T, error) {

	var data T

	if err := ctx.Bind(&data); err != nil {
		Logger.Error(err)
		return nil, err
	}

	if err := validatorInstance.Struct(data); err != nil {
		Logger.Error(err)
		return nil, err
	}

	return &data, nil
}

// validateEmoji validates whether the field is a valid emoji
func validateEmoji(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return false
	}

	// For emoji validation, we need to handle sequences properly
	// Check if the entire string consists of valid emoji characters and modifiers
	runes := []rune(value)
	if len(runes) == 0 {
		return false
	}

	// Allow emoji sequences with base emojis, modifiers, and joiners
	for _, r := range runes {
		if !isEmoji(r) {
			return false
		}
	}

	// Special case: allow zero-width joiner sequences
	// But ensure we have at least one base emoji
	hasBaseEmoji := false
	for _, r := range runes {
		// Base emoji ranges (excluding modifiers and joiners)
		if (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
			(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
			(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
			(r >= 0x1F1E6 && r <= 0x1F1FF) || // Regional country flags
			(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
			(r >= 0x2700 && r <= 0x27BF) || // Dingbats
			(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
			(r >= 0x1F018 && r <= 0x1F270) || // Various asian characters
			(r >= 0x238C && r <= 0x2454) { // Dingbats and additional emoticons
			hasBaseEmoji = true
			break
		}
	}

	// Must have at least one base emoji
	if !hasBaseEmoji {
		return false
	}

	// Limit emoji length (typically emojis are 1-4 characters, but sequences can be longer)
	return len(value) >= 1 && len(value) <= 20
}

// isEmoji checks if a single character is an emoji
func isEmoji(r rune) bool {
	// 基本 emoji 范围
	if (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E6 && r <= 0x1F1FF) || // Regional country flags
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1F018 && r <= 0x1F270) || // Various asian characters
		(r >= 0x238C && r <= 0x2454) || // Dingbats and additional emoticons
		(r >= 0x20D0 && r <= 0x20FF) || // Combining Diacritical Marks for Symbols
		(r >= 0x1F3FB && r <= 0x1F3FF) || // Skin tone modifiers
		(r == 0x200D) { // Zero-width joiner
		return true
	}

	// 检查是否为零宽连接符序列的一部分
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Cf) {
		return true
	}

	return false
}

func ValidateAvatarFile(fileHeader *multipart.FileHeader) ([]byte, string, error) {
	const (
		MAX_AVATAR_SIZE_BYTES = 5 << 20
		MAX_AVATAR_WIDTH      = 1024
		MAX_AVATAR_HEIGHT     = 1024
	)
	if fileHeader == nil {
		return nil, "", errors.New("avatar file is required")
	}
	if fileHeader.Size > MAX_AVATAR_SIZE_BYTES {
		return nil, "", errors.New("Avatar file is too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", errors.New("Failed to open file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", errors.New("Failed to open file")
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", errors.New("Invalid avatar format")
	}
	if config.Width > MAX_AVATAR_WIDTH || config.Height > MAX_AVATAR_HEIGHT {
		return nil, "", errors.New("Invalid avatar dimensions")
	}
	switch strings.ToLower(format) {
	case "jpeg", "png", "gif":
		return data, format, nil
	default:
		return nil, "", errors.New("Invalid avatar format")
	}
	return data, format, nil
}
