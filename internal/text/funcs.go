package text

import "fmt"

func GetConvertToLang(lang string, key int, formats ...any) string {
	switch lang {
	case "en":
		return fmt.Sprintf(En[key], formats...)
	case "ru":
		return fmt.Sprintf(Ru[key], formats...)
	default:
		return fmt.Sprintf(En[key], formats...)
	}
}