package main

import (
	"fmt"
	"log"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := botapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}
	activeChallenges := make(map[int]bool)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	timeouts := make(chan int)

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			//TODO listen to the channel where challenges get completed
			if update.Message.NewChatMembers != nil {
				for _, user := range *update.Message.NewChatMembers {
					//TODO extract the hardcoded messages
					m := fmt.Sprintf("Bem vinda ao DevOps Recife %s!\nEnviei um desafio para você no chat privado e espero que você me retorne em até 60 segundos ou terei que te convidar para sair do grupo! Nada pessoal, só não aceitamos spammers! :P", user.UserName)
					msg := botapi.NewMessage(int64(user.ID), m)
					bot.Send(msg)

					challenge := generateChallenge()
					challengeMsg := fmt.Sprintf("Quanto é %d %s %d? Você tem 60 segundos!", challenge.ElementA, challenge.Operation, challenge.ElementB)
					msg = botapi.NewMessage(int64(user.ID), challengeMsg)
					bot.Send(msg)

					activeChallenges[user.ID] = true
					go timeCounter(user.ID, timeouts)
				}

				// msg.ReplyToMessageID = update.Message.MessageID
			}
			if update.Message.Chat.IsPrivate() {
				//TODO check if there are challenges currently active for the given user and retrieve them
				fmt.Println("")
				var msg botapi.MessageConfig
				//TODO validate the answer according to the challenge
				if update.Message.Text == "8" {
					msg = botapi.NewMessage(update.Message.Chat.ID, "Resposta correta! Obrigado!")
					//TODO remove user from active challenge list
				} else {
					msg = botapi.NewMessage(update.Message.Chat.ID, "Não foi dessa vez, quer tentar de novo?")
				}
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}

			//TODO wipe all messages if a user post in the group and has an active challenge.
		case userID := <-timeouts:
			if activeChallenges[userID] {
				//TODO send message too
				fmt.Println("Kick the user!")
				delete(activeChallenges, userID)
			}
		}
	}
}

func timeCounter(userID int, channel chan<- int) {
	time.Sleep(60 * time.Second)
	channel <- userID
}
