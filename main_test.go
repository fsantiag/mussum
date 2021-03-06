package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/fsantiag/mussum/challenge"
	"github.com/fsantiag/mussum/language"
	"github.com/fsantiag/mussum/mocks"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendChallengeToUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := challenge.Generate()
	l := language.GetDefault()
	userName := "someUser"

	msg1 := botapi.NewMessage(20, fmt.Sprintf(l.Welcome(), userName))
	msg2 := botapi.NewMessage(20, fmt.Sprintf(l.Challenge(), c.A, c.Operation, c.B))

	mockBotAPI.EXPECT().Send(msg1).Return(botapi.Message{}, nil).Times(1)
	mockBotAPI.EXPECT().Send(msg2).Return(botapi.Message{}, nil).Times(1)
	sendChallengeToUser(userName, 20, l, mockBotAPI, c)
}

func TestUserPassesTheChallenge(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := challenge.Generate()
	u := botapi.Update{
		Message: &botapi.Message{
			Chat: &botapi.Chat{
				ID: 1,
			},
			Text: strconv.Itoa(c.Result),
			From: &botapi.User{
				ID: 1,
			},
			MessageID: 10,
		},
	}
	m := map[int]challenge.SumChallenge{
		1: c,
	}
	l := language.GetDefault()

	msg := botapi.NewMessage(u.Message.Chat.ID, l.Correct())
	msg.ReplyToMessageID = u.Message.MessageID

	mockBotAPI.EXPECT().Send(msg).Return(botapi.Message{}, nil).Times(1)
	verifyUserAnswer(u, m, l, mockBotAPI)
	assert.Contains(t, buf.String(), "[1] User successfully solved the challange")
	_, ok := m[u.Message.From.ID]
	assert.False(t, ok)
}

func TestUserFailsTheChallenge(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := challenge.Generate()
	u := botapi.Update{
		Message: &botapi.Message{
			Chat: &botapi.Chat{
				ID: 1,
			},
			Text: "wrong answer",
			From: &botapi.User{
				ID: 1,
			},
			MessageID: 10,
		},
	}
	m := map[int]challenge.SumChallenge{
		1: c,
	}
	l := language.GetDefault()

	msg := botapi.NewMessage(u.Message.Chat.ID, l.Wrong())
	msg.ReplyToMessageID = u.Message.MessageID

	mockBotAPI.EXPECT().Send(msg).Return(botapi.Message{}, nil).Times(1)
	verifyUserAnswer(u, m, l, mockBotAPI)
	assert.Contains(t, buf.String(), "[1] Wrong answer for challenge")
	_, ok := m[u.Message.From.ID]
	assert.True(t, ok)
}

// func TestIgnorePrivateMessagesWhenNoChallengeIsActiveForUser(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
// 	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
// 	var buf bytes.Buffer
// 	log.SetOutput(&buf)
// 	defer func() {
// 		log.SetOutput(os.Stdout)
// 	}()
// 	m := map[int]challenge.SumChallenge{}
// 	l := language.Pt{}
// 	u := botapi.Update{
// 		Message: &botapi.Message{
// 			Chat: &botapi.Chat{
// 				Type: "private",
// 			},
// 			From: &botapi.User{
// 				ID: 10,
// 			},
// 		},
// 	}
// 	updates := make(chan botapi.Update, 1)
// 	updates <- u
// 	quit := make(chan bool, 1)
// 	timeout := make(chan message, 1)
// 	quit <- true

// 	mockBotAPI.EXPECT().Send(gomock.Any()).Return(botapi.Message{}, nil).Times(0)
// 	startBot(mockBotAPI, l, m, updates, quit, timeout)
// }

func TestKickUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	kickCfg := botapi.KickChatMemberConfig{
		ChatMemberConfig: botapi.ChatMemberConfig{
			UserID: 555,
			ChatID: 2,
		},
		UntilDate: 400, //forever
	}

	mockBotAPI.EXPECT().KickChatMember(kickCfg).Return(botapi.APIResponse{}, nil).Times(1)

	kickUser(message{userID: 555, chatID: 2}, mockBotAPI)
	assert.Contains(t, buf.String(), "[555] User kicked from the channel")
}

func TestKickUserFails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotAPI := mocks.NewMockBotIface(mockCtrl)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	kickCfg := botapi.KickChatMemberConfig{
		ChatMemberConfig: botapi.ChatMemberConfig{
			UserID: 1,
			ChatID: 2,
		},
		UntilDate: 400, //forever
	}

	mockBotAPI.EXPECT().KickChatMember(kickCfg).Return(botapi.APIResponse{}, fmt.Errorf("something went wrong")).Times(1)

	kickUser(message{userID: 1, chatID: 2}, mockBotAPI)
	assert.Contains(t, buf.String(), "[1] Unable to kick user from group: something went wrong")
}
