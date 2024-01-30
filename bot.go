// main.go
package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/aswinxd/golang/config.go" // Import the config package
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
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
		bio, err := getUserBio(bot, user)
		if err != nil {
			log.Println("Error retrieving user bio:", err)
			continue
		}

		log.Printf("User %s bio: %s", user.UserName, bio)

		if bio != "" && hasLink(bio) {
			log.Printf("Banning user %s with link in bio", user.UserName)

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
			} else {
				log.Printf("User %s banned successfully", user.UserName)
			}
		}
	}
}

func getUserBio(bot *tgbotapi.BotAPI, user *tgbotapi.User) (string, error) {
	// Check if the latest message is present
	if update.Message.Text != "" {
		// Extract bio information from the latest message text
		return update.Message.Text, nil
	}

	// Return an empty string if no bio is available
	return "", nil
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
