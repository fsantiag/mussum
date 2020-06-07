package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fsantiag/mussum/adapter"
	"github.com/fsantiag/mussum/challenge"
	"github.com/fsantiag/mussum/language"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
	userID int
	chatID int64
}

func main() {
	bot, err := adapter.NewBotAPI(os.Getenv("APIKEY"))
	if err != nil {
		log.Fatalf("Unable to connect to Telegram Bot API: %v", err)
	}
	// bot.Debug = true
	lang := language.GetDefault()
	activeChallenges := make(map[int]challenge.SumChallenge)

	log.Printf("Bot started and authorized on account [%v]", bot.UserName())

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
					log.Printf("[%v] New user joined group", user.ID)

					c := challenge.Generate()

					sendChallengeToUser(user, lang, bot, c)

					activeChallenges[user.ID] = c

					log.Printf("[%v] Challenge sent for user", user.ID)

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
					verifyUserAnswer(update, challenge, activeChallenges, lang, bot)
				}
			}

			//TODO wipe all messages if a user post in the group and has an active challenge.
		case m := <-timeoutChannel:
			if _, ok := activeChallenges[m.userID]; ok {
				log.Printf("[%v] User failed to solve challenge", m.userID)
				kickUser(m, bot)
				delete(activeChallenges, m.userID)
			}
		}
	}
}

func sendChallengeToUser(
	user botapi.User,
	lang language.Language,
	bot adapter.BotIface,
	c challenge.SumChallenge) {

	msg := botapi.NewMessage(int64(user.ID), lang.Welcome())
	bot.Send(msg)

	msg = botapi.NewMessage(int64(user.ID), fmt.Sprintf(lang.Challenge(), c.A, c.Operation, c.B))
	bot.Send(msg)
}

func verifyUserAnswer(
	u botapi.Update,
	c challenge.SumChallenge,
	activeChallenges map[int]challenge.SumChallenge,
	lang language.Language,
	bot adapter.BotIface) {

	var msg botapi.MessageConfig
	if u.Message.Text == strconv.Itoa(c.Result) {
		msg = botapi.NewMessage(u.Message.Chat.ID, lang.Correct())
		delete(activeChallenges, u.Message.From.ID)

		log.Printf("[%v] User successfully solved the challange", u.Message.From.ID)
	} else {
		msg = botapi.NewMessage(u.Message.Chat.ID, lang.Wrong())

		log.Printf("[%v] Wrong answer for challenge", u.Message.From.ID)
	}
	msg.ReplyToMessageID = u.Message.MessageID
	bot.Send(msg)
}

func kickUser(message Message, bot adapter.BotIface) {

	kickCfg := botapi.KickChatMemberConfig{
		ChatMemberConfig: botapi.ChatMemberConfig{
			UserID: message.userID,
			ChatID: message.chatID,
		},
		UntilDate: 400, //forever
	}

	_, err := bot.KickChatMember(kickCfg)
	if err != nil {
		log.Printf("[%v] Unable to kick user from group: %v", message.userID, err)
	} else {
		log.Printf("[%v] User kicked from the channel", message.userID)
	}
}
