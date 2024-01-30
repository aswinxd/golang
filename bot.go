// main.go
package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const BotToken = "6350652136:AAFriDrVaXsIEchvLTj8BY3JEvvCGyVjTHI"

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
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
		bio, err := getUserBio(update.Message)
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

func getUserBio(message *tgbotapi.Message) (string, error) {
	// Check if the message text is present
	if message.Text != "" {
		// Extract bio information from the message text
		return message.Text, nil
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
