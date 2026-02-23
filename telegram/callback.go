package telegram

import (
	"GLPITGBOT/models"
	"GLPITGBOT/repository"
	"GLPITGBOT/telegram/i18n"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	q := update.CallbackQuery
	if q == nil || q.Message == nil {
		return
	}

	user, err := repository.EnsureUser(
		q.From.ID,
		q.From.UserName,
	)
	if err != nil {
		log.Println("ensure user error:", err)
		return
	}

	if user.Lang == "" {
		user.Lang = i18n.DetectLang(q.From.LanguageCode)
	}

	msg := tgbotapi.NewMessage(q.Message.Chat.ID, "")
	session := getSession(q.From.ID)

	switch q.Data {
	case CbCreateTicket:
		session = sessionData{
			Step:         StepWaitTicketTitle,
			ActiveTicket: &models.Ticket{},
		}
		setSession(q.From.ID, session)
		msg.Text = i18n.T(user.Lang, "enter_ticket_title")
		msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)

	case CbCancelTicket:
		resetSession(q.From.ID)
		msg.Text = i18n.T(user.Lang, "exit")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

	case CbLastTickets:
		result, err := repository.GetLastUserTicketsText(q.From.ID, 5)
		if err != nil {
			msg.Text = i18n.T(user.Lang, "error_get_tickets")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			break
		}

		msg.Text = result
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

	case CbPreferences:
		msg.Text = i18n.T(user.Lang, "preferences")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	case CbNotifications:
		msg.Text = i18n.T(user.Lang, "notifications")
		msg.ReplyMarkup = NotificationsKeyboard(user.Lang)

	case CbOpenLanguage:
		msg.Text = i18n.T(user.Lang, "choose_language")
		msg.ReplyMarkup = LanguageKeyboard()

	case CbLangRU:
		user.Lang = "ru"
		msg.Text = i18n.T(user.Lang, "language_changed")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	case CbLangEN:
		user.Lang = "en"
		msg.Text = i18n.T(user.Lang, "language_changed")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	case CbLangKK:
		user.Lang = "kk"
		msg.Text = i18n.T(user.Lang, "language_changed")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	case CbStart, CbExit, CbMainMenu:
		resetSession(q.From.ID)
		msg.Text = i18n.T(user.Lang, "start")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

	default:
		return
	}

	if session.Step == StepWaitTicketDesc && q.Data != CbCreateTicket {
		setSession(q.From.ID, session)
	}

	if err := repository.SaveUser(user); err != nil {
		log.Println("save user error:", err)
	}

	if _, err := bot.Send(tgbotapi.NewCallback(q.ID, "")); err != nil {
		log.Println("callback ack error:", err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println("send msg error:", err)
	}
}
