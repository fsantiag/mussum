package main

import (
	"fmt"
	"log"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := botapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
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

				//TODO generate a random challenge
				msg = botapi.NewMessage(int64(user.ID), fmt.Sprintf("Quanto é 4 + 4? Você tem 60 segundos: %s", user.UserName))
				bot.Send(msg)
				//TODO add user to the active challenge list
				//TODO start a goroutine to count the time and check if user passes the test
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
	}
}
