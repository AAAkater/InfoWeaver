package utils

import (
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

	// Check if it contains only emoji characters
	for _, r := range value {
		if !isEmoji(r) {
			return false
		}
	}

	// Limit emoji length (typically emojis are 1-4 characters)
	return len(value) >= 1 && len(value) <= 10
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
		(r >= 0x20D0 && r <= 0x20FF) { // Combining Diacritical Marks for Symbols
		return true
	}

	// 检查是否为零宽连接符序列的一部分
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Cf) {
		return true
	}

	return false
}
