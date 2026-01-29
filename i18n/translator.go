package i18n

import "GLPITGBOT/telegram/i18n"

func T(lang, key string) string {
	switch lang {
	case "ru":
		return i18n.RUS[key]
	case "kk":
		return i18n.KAZ[key]
	default:
		return i18n.ENG[key]
	}
}
func DetectLang(code string) string {
	switch code {
	case "ru", "en", "kk":
		return code
	default:
		return "en"
	}
}
