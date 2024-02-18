package commands

import (
  "fmt"
	. "collect3/renterd-telegram-alerts/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Unsubscribe (bot *tgbotapi.BotAPI, message *tgbotapi.MessageConfig) error {
  updated, err := DB.RemoveSubscription(message.ChatID)
  
  if err != nil {
    message.Text = "Failed to unsubscribe"
    fmt.Println(err.Error())
    _, err = bot.Send(message)
    return err
  }

  if !updated {
    message.Text = "You are not subscribed"
    _, err = bot.Send(message)
    return err
  }

  message.Text = "Unsubscribed succesfully"
  _, err = bot.Send(message)

  return err
}
