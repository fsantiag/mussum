package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := botapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}
	userActiveChallenges := make(map[int64]SumChallenge)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	timeoutChannel := make(chan int64)

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			//TODO listen to the channel where challenges get completed
			if update.Message.NewChatMembers != nil {
				for _, user := range *update.Message.NewChatMembers {
					userID := int64(user.ID)
					//TODO extract the hardcoded messages
					m := fmt.Sprintf("Bem vinda ao DevOps Recife %s!\nEnviei um desafio para você no chat privado e espero que você me retorne em até 60 segundos ou terei que te convidar para sair do grupo! Nada pessoal, só não aceitamos spammers! :P", user.UserName)
					msg := botapi.NewMessage(userID, m)
					bot.Send(msg)

					challenge := GenerateChallenge()
					challengeMsg := fmt.Sprintf("Quanto é %d %s %d? Você tem 60 segundos!", challenge.ElementA, challenge.Operation, challenge.ElementB)
					msg = botapi.NewMessage(userID, challengeMsg)
					bot.Send(msg)

					userActiveChallenges[userID] = challenge
					go timeCounter(userID, timeoutChannel)
				}
			}
			if update.Message.Chat.IsPrivate() {
				userID := update.Message.Chat.ID
				if challenge, ok := userActiveChallenges[userID]; ok {
					var msg botapi.MessageConfig
					if update.Message.Text == strconv.Itoa(challenge.Answer) {
						msg = botapi.NewMessage(userID, "Resposta correta! Obrigado!")
						delete(userActiveChallenges, userID)
					} else {
						msg = botapi.NewMessage(userID, "Não foi dessa vez, quer tentar de novo?")
					}
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}

			//TODO wipe all messages if a user post in the group and has an active challenge.
		case userID := <-timeoutChannel:
			if _, ok := userActiveChallenges[userID]; ok {
				//TODO kick the user from the group
				delete(userActiveChallenges, userID)
			}
		}
	}
}

func timeCounter(userID int64, channel chan<- int64) {
	time.Sleep(60 * time.Second)
	channel <- userID
}
