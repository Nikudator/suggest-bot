package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

func main() {
	//Читаем конфиг
	const configPath = "config.yml"
	type Cfg struct {
		TELEGRAM_BOT_API_TOKEN string `yaml:"token"`
	}
	var AppConfig *Cfg
	f, err := os.Open(configPath)

	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&AppConfig)

	if err != nil {
		log.Panic(err)
	}

	bot_token := AppConfig.TELEGRAM_BOT_API_TOKEN

	//Создаём бота
	bot, err := tgbotapi.NewBotAPI(bot_token)
	//
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Бот подключился %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var msg tgbotapi.MessageConfig

			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! I am your bot.")
			case "help":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I can help you with the following commands:\n/start - Start the bot\n/help - Display this help message")
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
			}

			bot.Send(msg)
		}
	}
}
