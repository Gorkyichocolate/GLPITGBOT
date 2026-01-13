package main

import (
	db "GLPITGBOT/DB"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создание Заявки"),
		tgbotapi.NewKeyboardButton("Мои Заявки"),
	),
)

//var ticketCreate = tgbotapi.

func main() {

	//Можно не трогать

	//--------------------БД--------------
	godotenv.Load()

	// DB
	if err := db.Connect(); err != nil {
		log.Fatal("db error:", err)
	}
	defer db.Pool.Close()

	//---------------------
	//----------Бот----------
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	//-------------------------------------------------

	for update := range updates {

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Command() {
			case "/start":
				msg.ReplyMarkup = startKeyboard
			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		}
	}
}
