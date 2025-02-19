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
		ADMIN_ID               int    `yaml:"admin_id"`
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
	admin_id := AppConfig.ADMIN_ID

	//Создаём бота
	bot, err := tgbotapi.NewBotAPI(bot_token)

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
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую! Я бот для отправки сообщений и постов на канал \"Реальное Шушенское\" t.me/real_shush \n Просто напишите сообщение мне и я передам его администратору")
			case "help":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Я поддерживаю следующие комманды:\n/start - Старт бота\n /help - Показать помощь\nЕсли хотите опубликовать пост или написать администратору сообщение, просто напишите его и, если нужно, прикрепите фото или видео.")
			default:
				var msg_adm tgbotapi.ForwardConfig
				msg_adm = tgbotapi.NewForward(int64(admin_id), update.Message.From.ID, update.Message.MessageID)
				bot.Send(msg_adm)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше сообщение отправлено администратору канала.")
			}

			bot.Send(msg)
		}
	}
}
