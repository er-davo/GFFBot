package text

import (
	"fmt"

	"gffbot/internal/logger"

	"go.uber.org/zap"
)

func init() {
	logger.Init()
}

func Convert(lang string, key int, formats ...any) string {
	logger.Log.Debug("text.Convert(...) called", zap.String("lang", lang))
	
	switch lang {
	case "en":
		return fmt.Sprintf(En[key], formats...)
	case "ru":
		return fmt.Sprintf(Ru[key], formats...)
	default:
		return fmt.Sprintf(En[key], formats...)
	}
}
