package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fsantiag/mussum/challenge"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
	userID int
	chatID int64
}

func main() {
	bot, err := botapi.NewBotAPI("")
	if err != nil {
		log.Fatalf("Unable to connect to Telegram Bot API: %v", err)
	}
	activeChallenges := make(map[int]challenge.SumChallenge)

	bot.Debug = true

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
					//TODO extract the hardcoded messages
					m := fmt.Sprintf("Bem vinda ao DevOps Recife %s!\nEnviei um desafio para você no chat privado e espero que você me retorne em até 60 segundos ou terei que te convidar para sair do grupo! Nada pessoal, só não aceitamos spammers! :P", user.UserName)
					msg := botapi.NewMessage(int64(user.ID), m)
					bot.Send(msg)

					c := challenge.Generate()
					challengeMsg := fmt.Sprintf("Quanto é %d %s %d? Você tem 60 segundos!", c.ElementA, c.Operation, c.ElementB)
					msg = botapi.NewMessage(int64(user.ID), challengeMsg)
					bot.Send(msg)

					activeChallenges[user.ID] = c

					go timer(Message{
						userID: user.ID,
						chatID: update.Message.Chat.ID,
					}, timeoutChannel)
				}
			}
			if update.Message.Chat.IsPrivate() {
				if challenge, ok := activeChallenges[update.Message.From.ID]; ok {
					var msg botapi.MessageConfig
					if update.Message.Text == strconv.Itoa(challenge.Answer) {
						msg = botapi.NewMessage(update.Message.Chat.ID, "Resposta correta! Obrigado!")
						delete(activeChallenges, update.Message.From.ID)
					} else {
						msg = botapi.NewMessage(update.Message.Chat.ID, "Não foi dessa vez, quer tentar de novo?")
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

func timer(message Message, channel chan<- Message) {
	time.Sleep(60 * time.Second)
	channel <- message
}
