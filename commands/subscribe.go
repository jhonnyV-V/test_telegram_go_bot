package commands

import (
	. "collect3/renterd-telegram-alerts/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Subscribe (bot *tgbotapi.BotAPI, message *tgbotapi.MessageConfig) error {
  updated, err := DB.AddSubscription(message.ChatID)
  
  if err != nil {
    message.Text = "Failed to subscribe"
    _, err = bot.Send(message)
    return err
  }
  if !updated {
    message.Text = "You are already subscribed"
    _, err = bot.Send(message)
    return err
  }

  message.Text = "Subscribed succesfully"
  _, err = bot.Send(message)

  return err
}
