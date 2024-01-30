package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1711796263:AAGwSl_kQXts-Q4Q_5NjuWcgUcncinO8M7M")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		user := update.Message.From
		chatID := update.Message.Chat.ID
		bio := user.Bio

		if bio != "" && hasLink(bio) {
			kickConfig := tgbotapi.KickChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: chatID,
					UserID: user.ID,
				},
				UntilDate: 0, // Ban permanently
			}

			_, err := bot.KickChatMember(kickConfig)
			if err != nil {
				log.Println("Failed to ban user:", err)
			}
		}
	}
}

func hasLink(text string) bool {
	patterns := []string{
		`@`,
		`https://`,
		`http://`,
		`t\.me//`,
		`t\.me`,
	}

	for _, pattern := range patterns {
		match, _ := regexp.MatchString(pattern, text)
		if match {
			return true
		}
	}

	return false
}
