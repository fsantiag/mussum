package adapter

import botapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BotIface interface {
	KickChatMember(config botapi.KickChatMemberConfig) (botapi.APIResponse, error)
	GetUpdatesChan(config botapi.UpdateConfig) (botapi.UpdatesChannel, error)
	Send(c botapi.Chattable) (botapi.Message, error)
	DeleteMessage(config botapi.DeleteMessageConfig) (botapi.APIResponse, error)
	UserName() string
}

type BotAdapter struct {
	Bot  *botapi.BotAPI
	Self botapi.User
}

func NewBotAPI(token string, debug bool) (BotIface, error) {
	b, err := botapi.NewBotAPI(token)
	b.Debug = debug
	if err != nil {
		return nil, err
	}
	return BotAdapter{b, b.Self}, nil
}

func (b BotAdapter) KickChatMember(config botapi.KickChatMemberConfig) (botapi.APIResponse, error) {
	return b.Bot.KickChatMember(config)
}
func (b BotAdapter) GetUpdatesChan(config botapi.UpdateConfig) (botapi.UpdatesChannel, error) {
	return b.Bot.GetUpdatesChan(config)
}
func (b BotAdapter) Send(c botapi.Chattable) (botapi.Message, error) {
	return b.Bot.Send(c)
}

func (b BotAdapter) DeleteMessage(config botapi.DeleteMessageConfig) (botapi.APIResponse, error) {
	return b.Bot.DeleteMessage(config)
}
func (b BotAdapter) UserName() string {
	return b.Bot.Self.UserName
}
