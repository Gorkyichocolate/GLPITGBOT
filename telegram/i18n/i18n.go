package i18n

var ENG = map[string]string{
	"start":            "Welcome!",
	"creating_tickets": "Create Ticket",
	"last_tickets":     "Last Tickets",
	"notifications":    "Notifications",
	"preferences":      "Preferences",
	"language":         "Language",
	"exit":             "Exit",

	"btn_exit":          "Exit",
	"btn_start":         "Start",
	"btn_language":      "Change language",
	"btn_notifications": "Notifications",

	"enter_ticket_title":       "Enter ticket title",
	"enter_ticket_description": "Enter ticket description",
	"ticket_created":           "✅ Ticket created successfully",

	"error_create_ticket": "❌ Error creating ticket",
	"error_get_tickets":   "❌ Error getting tickets",

	"choose_menu":      "Please choose an option from the menu",
	"choose_language":  "Choose language",
	"language_changed": "Language changed successfully",
	"enter_api_key":    "Enter your API key",
}

var RUS = map[string]string{
	"start":            "Добро пожаловать!",
	"creating_tickets": "Создать заявку",
	"last_tickets":     "Последние заявки",
	"notifications":    "Уведомления",
	"preferences":      "Настройки",
	"language":         "Язык",
	"exit":             "Выход",

	"btn_exit":          "Выйти",
	"btn_start":         "Старт",
	"btn_language":      "Сменить язык",
	"btn_notifications": "Уведомления",

	"enter_ticket_title":       "Введите тему заявки",
	"enter_ticket_description": "Введите описание заявки",
	"ticket_created":           "✅ Заявка успешно создана",

	"error_create_ticket": "❌ Ошибка при создании заявки",
	"error_get_tickets":   "❌ Ошибка при получении заявок",

	"choose_menu":      "Пожалуйста, выберите пункт меню",
	"choose_language":  "Выберите язык",
	"language_changed": "Язык успешно изменён",
	"enter_api_key":    "Введите API ключ",
}

var KAZ = map[string]string{
	"start":            "Қош келдіңіз!",
	"creating_tickets": "Өтінім жасау",
	"last_tickets":     "Соңғы өтінімдер",
	"notifications":    "Хабарламалар",
	"preferences":      "Баптаулар",
	"language":         "Тіл",
	"exit":             "Шығу",

	"btn_exit":          "Шығу",
	"btn_start":         "Бастау",
	"btn_language":      "Тілді ауыстыру",
	"btn_notifications": "Хабарламалар",

	"enter_ticket_title":       "Өтінім тақырыбын енгізіңіз",
	"enter_ticket_description": "Өтінім сипаттамасын енгізіңіз",
	"ticket_created":           "✅ Өтінім сәтті жасалды",

	"error_create_ticket": "❌ Өтінім жасау кезінде қате",
	"error_get_tickets":   "❌ Өтінімдерді алу қатесі",

	"choose_menu":      "Мәзірден таңдаңыз",
	"choose_language":  "Тілді таңдаңыз",
	"language_changed": "Тіл сәтті өзгертілді",
	"enter_api_key":    "API кілтін енгізіңіз",
}

var translations = map[string]map[string]string{
	"en": ENG,
	"ru": RUS,
	"kk": KAZ,
}

func T(lang, key string) string {
	if langMap, ok := translations[lang]; ok {
		if val, ok := langMap[key]; ok {
			return val
		}
	}

	if val, ok := RUS[key]; ok {
		return val
	}

	return key
}

func DetectLang(code string) string {
	if len(code) < 2 {
		return "ru"
	}

	switch code[:2] {
	case "ru":
		return "ru"
	case "kk":
		return "kk"
	case "en":
		return "en"
	default:
		return "ru"
	}
}
