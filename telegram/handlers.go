package telegram

import (
	glpihttp "GLPITGBOT/https"
	"GLPITGBOT/models"
	"GLPITGBOT/repository"
	"GLPITGBOT/telegram/i18n"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dueDateLayout = "2006-01-02 15:04:05"

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
	text := strings.TrimSpace(update.Message.Text)

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

	if session.Step == StepWaitApiToken {
		if text == "" {
			session.Step = StepWaitApiToken
			setSession(telegramID, session)
			msg.Text = i18n.T(user.Lang, "need_auth_first") + "\n" + i18n.T(user.Lang, "enter_api_key")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			if _, sendErr := bot.Send(msg); sendErr != nil {
				log.Println("SEND ERROR:", sendErr)
			}
			return
		}

		sessionToken, err := glpihttp.AuthByUserToken(text)
		if err != nil {
			session.Step = StepWaitApiToken
			setSession(telegramID, session)
			msg.Text = i18n.T(user.Lang, "auth_failed") + "\n" + i18n.T(user.Lang, "enter_api_key")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			if _, sendErr := bot.Send(msg); sendErr != nil {
				log.Println("SEND ERROR:", sendErr)
			}
			return
		}

		user.ApiToken = text
		user.SessionToken = sessionToken
		if err := repository.SaveUser(user); err != nil {
			log.Println("save user error:", err)
			msg.Text = i18n.T(user.Lang, "auth_failed")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			if _, sendErr := bot.Send(msg); sendErr != nil {
				log.Println("SEND ERROR:", sendErr)
			}
			return
		}

		resetSession(telegramID)
		msg.Text = i18n.T(user.Lang, "auth_success")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
		if _, err := bot.Send(msg); err != nil {
			log.Println("SEND ERROR:", err)
		}
		return
	}

	if !repository.IsUserAuthorized(user) {
		session.Step = StepWaitApiToken
		setSession(telegramID, session)
		msg.Text = i18n.T(user.Lang, "need_auth_first") + "\n" + i18n.T(user.Lang, "enter_api_key")
		msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
		if _, err := bot.Send(msg); err != nil {
			log.Println("SEND ERROR:", err)
		}
		return
	}

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

		session.ActiveTicket.Description = text
		session.Step = StepWaitTicketPriority
		setSession(telegramID, session)

		msg.Text = i18n.T(user.Lang, "enter_ticket_priority")
		msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)

	case StepWaitTicketPriority:
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

		priority, err := strconv.Atoi(text)
		if err != nil || priority < 0 || priority > 5 {
			msg.Text = i18n.T(user.Lang, "invalid_ticket_priority")
			msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)
			break
		}

		session.ActiveTicket.Priority = priority
		session.Step = StepWaitTicketDueDate
		setSession(telegramID, session)

		msg.Text = i18n.T(user.Lang, "enter_ticket_due_date")
		msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)

	case StepWaitTicketDueDate:
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

		if _, err := time.Parse(dueDateLayout, text); err != nil {
			msg.Text = i18n.T(user.Lang, "invalid_ticket_due_date")
			msg.ReplyMarkup = CancelTicketKeyboard(user.Lang)
			break
		}

		input := models.CreateTicketInput{
			Name:       session.ActiveTicket.Title,
			Content:    session.ActiveTicket.Description,
			EntitiesID: 0,
			Status:     1,
			Priority:   session.ActiveTicket.Priority,
			DueDate:    text,
		}

		localTicketID, err := repository.CreateTicketByTelegramIDWithStatus(
			telegramID,
			session.ActiveTicket.Title,
			session.ActiveTicket.Description,
			repository.TicketStatusWaiting,
		)
		if err != nil {
			log.Println("save local ticket error:", err)
			msg.Text = i18n.T(user.Lang, "error_create_ticket")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			resetSession(telegramID)
			break
		}

		if _, _, err := glpihttp.CreateTicketWithSession(user.SessionToken, input); err != nil {
			if updErr := repository.UpdateTicketStatus(localTicketID, repository.TicketStatusWaiting); updErr != nil {
				log.Println("update local ticket status error:", updErr)
			}

			msg.Text = i18n.T(user.Lang, "error_create_ticket")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)
			resetSession(telegramID)
			break
		}

		if err := repository.UpdateTicketStatus(localTicketID, repository.TicketStatusReview); err != nil {
			log.Println("update local ticket status error:", err)
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

		case i18n.T(user.Lang, "btn_logout"):
			user.ApiToken = ""
			user.SessionToken = ""
			if err := repository.SaveUser(user); err != nil {
				log.Println("logout save user error:", err)
			}
			resetSession(telegramID)
			session = getSession(telegramID)
			session.Step = StepWaitApiToken
			setSession(telegramID, session)
			msg.Text = i18n.T(user.Lang, "logged_out") + "\n" + i18n.T(user.Lang, "need_auth_first") + "\n" + i18n.T(user.Lang, "enter_api_key")
			msg.ReplyMarkup = MainMenuKeyboard(user.Lang)

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
				log.Println("get last tickets error:", err)
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

	if err := repository.SaveUser(user); err != nil {
		log.Println("save user error:", err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println("SEND ERROR:", err)
	}
}
