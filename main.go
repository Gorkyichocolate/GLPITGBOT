package main

import (
	"GLPITGBOT/db"
	"GLPITGBOT/telegram"
)

func main() {
	db.Connect()
	telegram.Bot()

}
