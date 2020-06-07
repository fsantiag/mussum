package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/fsantiag/mussum/mocks"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestKickUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotApi := mocks.NewMockBotIface(mockCtrl)
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

	mockBotApi.EXPECT().KickChatMember(kickCfg).Return(botapi.APIResponse{}, nil).Times(1)

	kickUser(Message{userID: 555, chatID: 2}, mockBotApi)
	assert.Contains(t, buf.String(), "[555] User kicked from the channel")
}

func TestKickUserFails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBotApi := mocks.NewMockBotIface(mockCtrl)
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

	mockBotApi.EXPECT().KickChatMember(kickCfg).Return(botapi.APIResponse{}, fmt.Errorf("Something went wrong!")).Times(1)

	kickUser(Message{userID: 1, chatID: 2}, mockBotApi)
	assert.Contains(t, buf.String(), "Unable to kick user from group: Something went wrong!")
}
