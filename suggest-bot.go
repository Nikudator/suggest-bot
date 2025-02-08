package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

func main() {
	//Читаем конфиг
	const configPath = "config.yml"
	type Cfg struct {
		TELEGRAM_BOT_API_TOKEN string `yaml:"token"`
		POSTGRES_HOST          string `yaml:"postgres_host"`
		POSTGRES_PORT          string `yaml:"postgres_port"`
		POSTGRES_DB            string `yaml:"postgres_db"`
		POSTGRES_USER          string `yaml:"postgres_user"`
		POSTGRES_PASS          string `yaml:"postgres_pass"`
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
	postgres_host := AppConfig.POSTGRES_HOST
	postgres_port := AppConfig.POSTGRES_PORT
	postgres_db := AppConfig.POSTGRES_DB
	postgres_user := AppConfig.POSTGRES_USER
	postgres_pass := AppConfig.POSTGRES_PASS

	//Инициализация БД
	psqlInfo := fmt.Sprintf("postgres_host=%s postgres_port=%d postgres_user=%s "+
		"postgres_pass=%s postgres_db=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

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
