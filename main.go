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

type message struct {
	userID int
	chatID int64
}

func main() {
	b, err := adapter.NewBotAPI(os.Getenv("APIKEY"), false)
	if err != nil {
		log.Fatalf("Unable to connect to Telegram Bot API: %v", err)
	}
	log.Printf("Bot started and authorized on account [%v]", b.UserName())
	m := make(map[int]challenge.SumChallenge)
	l := language.GetDefault()
	log.Printf("Active language: %v", l.Id())
	u := botapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Failed to retrieve update channel")
	}
	timeout := make(chan message)
	startBot(b, l, m, updates, timeout)
}

func startBot(
	bot adapter.BotIface,
	lang language.Language,
	activeChallenges map[int]challenge.SumChallenge,
	updates botapi.UpdatesChannel,
	timeout chan message) {

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			if update.Message.NewChatMembers != nil {
				for _, user := range *update.Message.NewChatMembers {
					if user.UserName == bot.UserName() {
						log.Printf("Bot joined group: [%v]", update.Message.Chat.Title)
						continue
					}
					log.Printf("[%v] New user joined group", user.ID)

					c := challenge.Generate()

					sendChallengeToUser(message{userID: user.ID, chatID: update.Message.Chat.ID}, lang, bot, c)

					activeChallenges[user.ID] = c

					log.Printf("[%v] Challenge sent for user", user.ID)

					go func() {
						time.Sleep(60 * time.Second)
						timeout <- message{
							userID: user.ID,
							chatID: update.Message.Chat.ID,
						}
					}()

				}
			}
			if !update.Message.Chat.IsPrivate() && update.Message.Text != "" {
				verifyUserAnswer(update, activeChallenges, lang, bot)
			}
			// else {
			// 	queue = append(queue, update)
			// 	bot.DeleteMessage(botapi.DeleteMessageConfig{
			// 		ChatID:    update.Message.Chat.ID,
			// 		MessageID: update.Message.MessageID,
			// 	})
			// }
		case m := <-timeout:
			if _, ok := activeChallenges[m.userID]; ok {
				log.Printf("[%v] User failed to solve challenge", m.userID)
				kickUser(m, bot)
				delete(activeChallenges, m.userID)
				continue
			}
			log.Printf("[%v] Timeout reached and user succeeded in challenge", m.userID)
		}
	}
}

func sendChallengeToUser(
	message message,
	lang language.Language,
	bot adapter.BotIface,
	c challenge.SumChallenge) {

	msg := botapi.NewMessage(message.chatID, lang.Welcome())
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("[%v] error: %v", message.chatID, err)
	}

	msg = botapi.NewMessage(message.chatID, fmt.Sprintf(lang.Challenge(), c.A, c.Operation, c.B))
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("[%v] error: %v", message.chatID, err)
	}
}

func verifyUserAnswer(
	u botapi.Update,
	activeChallenges map[int]challenge.SumChallenge,
	lang language.Language,
	bot adapter.BotIface) {

	if c, ok := activeChallenges[u.Message.From.ID]; ok {
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
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("[%v] error: %v", u.Message.From.ID, err)
		}
	}
}

func kickUser(message message, bot adapter.BotIface) {

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
