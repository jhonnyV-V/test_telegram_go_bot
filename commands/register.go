package commands

import (
  "fmt"
  "strings"
  "net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Register (
  bot *tgbotapi.BotAPI,
  message *tgbotapi.MessageConfig,
  command string,
  rawMessage string,
) error {
  var err error
  commandWithSlash := fmt.Sprintf("/%s", command)
  cleanMessage := strings.TrimSpace(
    strings.Replace(
      rawMessage,
      commandWithSlash,
      "",
      1,
    ),
  )
  args := strings.Split(cleanMessage, " ")
  numOfArgs := len(args)
  var webhook *url.URL
  module := "alerts"
  events := "register"

  switch numOfArgs {
    case 1:
      if args[0] == "" {
        message.Text = "You need to add a url to register"
        _, err = bot.Send(message)
        return err
      }
      webhook, err = url.ParseRequestURI(args[0])
      if err != nil || webhook.Scheme == "" || webhook.Host == "" {
        message.Text = "Invalid Url"
        fmt.Println("invalidurl: ", err.Error())
        _, err = bot.Send(message)
        return err
      }
      message.Text = fmt.Sprintf(
        "your endpoint %s",
        webhook.String(),
      )
      _, err = bot.Send(message)
    case 2:
      webhook, err = url.ParseRequestURI(args[0])
      module = args[1]
      if err != nil || webhook.Scheme == "" || webhook.Host == "" {
        message.Text = "Invalid Url"
        fmt.Println("invalidurl: ", err.Error())
        _, err = bot.Send(message)
        return err
      }
      message.Text = fmt.Sprintf(
        "your endpoint %s \nyour module %s",
        webhook.String(),
        module,
      )
      _, err = bot.Send(message)
    case 3:
      webhook, err = url.ParseRequestURI(args[0])
      module = args[1]
      events = args[2]
      if err != nil || webhook.Scheme == "" || webhook.Host == "" {
        message.Text = "Invalid Url"
        fmt.Println("invalidurl: ", err.Error())
        _, err = bot.Send(message)
        return err
      }
      message.Text = fmt.Sprintf(
        "your endpoint %s \nyour module %s\nyour event %s",
        webhook.String(),
        module,
        events,
      )
      _, err = bot.Send(message)

    default:
      message.Text = "Invalid Number of Arguments"
      _, err = bot.Send(message)
  }

  if err != nil {
    return err
  }
  //TODO: do the actual register of the webhook
  //TODO: save data in database
  return err
}
