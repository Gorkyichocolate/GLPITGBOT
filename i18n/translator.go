package i18n

import "GLPITGBOT/telegram/languages"

func T(lang, key string) string {
	switch lang {
	case "ru":
		return languages.RU[key]
	case "kk":
		return languages.KZ[key]
	default:
		return languages.ENG[key]
	}
}
