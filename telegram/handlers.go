package telegram

import (
	"GLPITGBOT/models"
	"GLPITGBOT/repository"
	"GLPITGBOT/telegram/i18n"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	if update.CallbackQuery != nil {
		handleCallback(bot, update)
		return
	}

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	text := update.Message.Text

	user, err := repository.EnsureUser(
		update.Message.From.ID,
		update.Message.From.UserName,
	)
	if err != nil {
		log.Println("ensure user error:", err)
		return
	}

	if user.Lang == "" {
		user.Lang = i18n.DetectLang(update.Message.From.LanguageCode)
		_ = repository.SaveUser(user)
	}

	msg := tgbotapi.NewMessage(chatID, "")
	session := getSession(telegramID)

	if text == "/start" {
		resetSession(telegramID)
		msg.Text = i18n.T(user.Lang, "start")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
		if _, err := bot.Send(msg); err != nil {
			log.Println("SEND ERROR:", err)
		}
		return
	}

	switch session.Step {
	case StepWaitTicketTitle:
		if session.ActiveTicket == nil {
			session.ActiveTicket = &models.Ticket{}
		}

		session.ActiveTicket.Title = text
		session.Step = StepWaitTicketDesc
		setSession(telegramID, session)

		msg.Text = i18n.T(user.Lang, "enter_ticket_description")
		msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)

	case StepWaitTicketDesc:
		if session.ActiveTicket == nil {
			session = sessionData{
				Step:         StepWaitTicketTitle,
				ActiveTicket: &models.Ticket{},
			}
			setSession(telegramID, session)
			msg.Text = i18n.T(user.Lang, "enter_ticket_title")
			msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)
			break
		}

		err := repository.CreateTicketByTelegramID(
			telegramID,
			session.ActiveTicket.Title,
			text,
		)

		if err != nil {
			msg.Text = i18n.T(user.Lang, "error_create_ticket")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			resetSession(telegramID)
			break
		}

		msg.Text = i18n.T(user.Lang, "ticket_created")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
		resetSession(telegramID)

	default:
		switch text {
		case i18n.T(user.Lang, "preferences"):
			msg.Text = i18n.T(user.Lang, "preferences")
			msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

		case i18n.T(user.Lang, "btn_notifications"):
			msg.Text = i18n.T(user.Lang, "notifications")
			msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

		case i18n.T(user.Lang, "btn_language"):
			msg.Text = i18n.T(user.Lang, "choose_language")
			msg.ReplyMarkup = LanguageKeyboard()

		case i18n.T(user.Lang, "btn_exit"):
			resetSession(telegramID)
			msg.Text = i18n.T(user.Lang, "start")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

		case i18n.T(user.Lang, "creating_tickets"):
			session = sessionData{
				Step:         StepWaitTicketTitle,
				ActiveTicket: &models.Ticket{},
			}
			setSession(telegramID, session)
			msg.Text = i18n.T(user.Lang, "enter_ticket_title")
			msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)

		case i18n.T(user.Lang, "last_tickets"):
			result, err := repository.GetLastUserTicketsText(telegramID, 5)
			if err != nil {
				msg.Text = i18n.T(user.Lang, "error_get_tickets")
				msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
				break
			}

			msg.Text = result
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

		default:
			msg.Text = i18n.T(user.Lang, "choose_menu")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
		}
	}

	_ = repository.SaveUser(user)

	if _, err := bot.Send(msg); err != nil {
		log.Println("SEND ERROR:", err)
	}
}
