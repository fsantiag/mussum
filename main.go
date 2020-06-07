package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fsantiag/mussum/challenge"
	"github.com/fsantiag/mussum/language"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
	userID int
	chatID int64
}

func main() {
	bot, err := botapi.NewBotAPI(os.Getenv("APIKEY"))
	lang := language.GetDefault()

	// bot.Debug = true
	if err != nil {
		log.Fatalf("Unable to connect to Telegram Bot API: %v", err)
	}
	activeChallenges := make(map[int]challenge.SumChallenge)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	timeoutChannel := make(chan Message)

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			if update.Message.NewChatMembers != nil {
				for _, user := range *update.Message.NewChatMembers {
					msg := botapi.NewMessage(int64(user.ID), lang.Welcome())
					bot.Send(msg)

					c := challenge.Generate()
					msg = botapi.NewMessage(int64(user.ID), fmt.Sprintf(lang.Challenge(), c.ElementA, c.Operation, c.ElementB))
					bot.Send(msg)

					activeChallenges[user.ID] = c

					go func() {
						time.Sleep(60 * time.Second)
						timeoutChannel <- Message{
							userID: user.ID,
							chatID: update.Message.Chat.ID,
						}
					}()
				}
			}
			if update.Message.Chat.IsPrivate() {
				if challenge, ok := activeChallenges[update.Message.From.ID]; ok {
					var msg botapi.MessageConfig
					if update.Message.Text == strconv.Itoa(challenge.Answer) {
						msg = botapi.NewMessage(update.Message.Chat.ID, lang.Correct())
						delete(activeChallenges, update.Message.From.ID)
					} else {
						msg = botapi.NewMessage(update.Message.Chat.ID, lang.Wrong())
					}
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}

			//TODO wipe all messages if a user post in the group and has an active challenge.
		case message := <-timeoutChannel:
			if _, ok := activeChallenges[message.userID]; ok {
				kickCfg := botapi.KickChatMemberConfig{
					ChatMemberConfig: botapi.ChatMemberConfig{
						UserID: message.userID,
						ChatID: message.chatID,
					},
					UntilDate: 400,
				}
				_, err := bot.KickChatMember(kickCfg)
				if err != nil {
					log.Printf("Unable to kick user from group: %v", err)
				}
				delete(activeChallenges, message.userID)
			}
		}
	}
}
