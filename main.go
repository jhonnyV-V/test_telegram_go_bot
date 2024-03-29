package main

import (
	"log"
  "os"
	. "collect3/renterd-telegram-alerts/utils"
	. "collect3/renterd-telegram-alerts/commands"

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
  db, err := OpenDB("sqlite3", "./db/local.db")
  if err != nil {
		log.Panic(err)
	}

  db.Migrate()

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

    command := update.Message.Command()

    //ignore most comments above, change in perspective
    //TODO: 5 commands
    //1- subscribe
    //2- unsubscribe
    //3- create_listener
    //4- listen
    //5- silence
    //maybe others like stats
    var err error
		switch command {
		case "help":
			msg.Text = "I understand /sayhi, /status and /register."
      _, err = bot.Send(msg)
		case "sayhi":
			msg.Text = "Hi :)"
      _, err = bot.Send(msg)
		case "status":
			msg.Text = "I'm ok."
      _, err = bot.Send(msg)
    case "subscribe":
      err = Subscribe(bot, &msg)
    case "unsubscribe":
      err = Unsubscribe(bot, &msg)
    //case "register": err = Register( bot, &msg, command, update.Message.Text,)
		default:
			msg.Text = "I don't know that command"
      _, err = bot.Send(msg)
		}
		if err != nil {
			log.Panic(err)
		}
	}
}
