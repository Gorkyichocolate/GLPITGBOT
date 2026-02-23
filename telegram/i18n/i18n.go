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
	"enter_ticket_priority":    "Enter priority from 0 to 5",
	"enter_ticket_due_date":    "Enter due date in format YYYY-MM-DD HH:MM:SS",
	"ticket_created":           "✅ Ticket created successfully",
	"auth_success":             "✅ Authorization successful",
	"auth_failed":              "❌ Authorization failed. Send valid user API token",

	"error_create_ticket":     "❌ Error creating ticket",
	"error_get_tickets":       "❌ Error getting tickets",
	"invalid_ticket_priority": "❌ Priority must be a number from 0 to 5",
	"invalid_ticket_due_date": "❌ Invalid date format. Use YYYY-MM-DD HH:MM:SS",

	"choose_menu":      "Please choose an option from the menu",
	"choose_language":  "Choose language",
	"language_changed": "Language changed successfully",
	"enter_api_key":    "Send your user API token to authorize",
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
	"enter_ticket_priority":    "Введите приоритет от 0 до 5",
	"enter_ticket_due_date":    "Введите дату дедлайна в формате ГГГГ-ММ-ДД ЧЧ:ММ:СС",
	"ticket_created":           "✅ Заявка успешно создана",
	"auth_success":             "✅ Авторизация прошла успешно",
	"auth_failed":              "❌ Ошибка авторизации. Отправьте корректный user API token",

	"error_create_ticket":     "❌ Ошибка при создании заявки",
	"error_get_tickets":       "❌ Ошибка при получении заявок",
	"invalid_ticket_priority": "❌ Приоритет должен быть числом от 0 до 5",
	"invalid_ticket_due_date": "❌ Неверный формат даты. Используйте ГГГГ-ММ-ДД ЧЧ:ММ:СС",

	"choose_menu":      "Пожалуйста, выберите пункт меню",
	"choose_language":  "Выберите язык",
	"language_changed": "Язык успешно изменён",
	"enter_api_key":    "Отправьте user API token для авторизации",
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
	"enter_ticket_priority":    "0-ден 5-ке дейін приоритет енгізіңіз",
	"enter_ticket_due_date":    "Дедлайнды YYYY-MM-DD HH:MM:SS форматында енгізіңіз",
	"ticket_created":           "✅ Өтінім сәтті жасалды",
	"auth_success":             "✅ Авторизация сәтті өтті",
	"auth_failed":              "❌ Авторизация қатесі. Дұрыс user API token жіберіңіз",

	"error_create_ticket":     "❌ Өтінім жасау кезінде қате",
	"error_get_tickets":       "❌ Өтінімдерді алу қатесі",
	"invalid_ticket_priority": "❌ Приоритет 0-ден 5-ке дейінгі сан болуы керек",
	"invalid_ticket_due_date": "❌ Күн форматы қате. YYYY-MM-DD HH:MM:SS қолданыңыз",

	"choose_menu":      "Мәзірден таңдаңыз",
	"choose_language":  "Тілді таңдаңыз",
	"language_changed": "Тіл сәтті өзгертілді",
	"enter_api_key":    "Авторизация үшін user API token жіберіңіз",
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
