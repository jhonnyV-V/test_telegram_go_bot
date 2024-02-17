package main

import (
	"log"
  "os"
  "fmt"
	//. "collect3/renterd-telegram-alerts/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func getEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
    log.Panic(err)
	}

	return os.Getenv(key)
}

func main() {
  //db, err := OpenDB("sqlite3", "./db/local.db")
  //if err != nil {
	//	log.Panic(err)
	//}

  //db.Migrate()

	bot, err := tgbotapi.NewBotAPI(getEnvVar("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

  if getEnvVar("ENV") != "prod" {
    bot.Debug = true
  }

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
    chatID := update.Message.Chat.ID
		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(chatID, "")
    //TODO: use the renterd /api/bus/webhooks/action to send custom alerts
    //TODO: use the renterd /api/bus/webhooks to register the webhooks
    //TODO: use the renter /api/bus/webhook/delete to delete webhooks

		// Extract the command from the Message.
    rawMessage := update.Message.Text
    command := update.Message.Command()
    fmt.Println(rawMessage)
    //TODO: take the message remove the command part, trim it,
    //check if it have a space, in that case split it and take the strings as arguments
    
    //TODO: to send the alerts check the origin url/ip of the caller of the endpoint 
		switch command {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
