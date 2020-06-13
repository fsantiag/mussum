package adapter

import botapi "github.com/go-telegram-bot-api/telegram-bot-api"

// BotIface is wrapper to the telegram-bot-api to allow testing the bot
type BotIface interface {
	KickChatMember(config botapi.KickChatMemberConfig) (botapi.APIResponse, error)
	GetUpdatesChan(config botapi.UpdateConfig) (botapi.UpdatesChannel, error)
	Send(c botapi.Chattable) (botapi.Message, error)
	DeleteMessage(config botapi.DeleteMessageConfig) (botapi.APIResponse, error)
	UserName() string
}

type botAdapter struct {
	Bot  *botapi.BotAPI
	Self botapi.User
}

// NewBotAPI creates a object to talk to the telegram API
func NewBotAPI(token string, debug bool) (BotIface, error) {
	b, err := botapi.NewBotAPI(token)
	b.Debug = debug
	if err != nil {
		return nil, err
	}
	return botAdapter{b, b.Self}, nil
}

func (b botAdapter) KickChatMember(config botapi.KickChatMemberConfig) (botapi.APIResponse, error) {
	return b.Bot.KickChatMember(config)
}
func (b botAdapter) GetUpdatesChan(config botapi.UpdateConfig) (botapi.UpdatesChannel, error) {
	return b.Bot.GetUpdatesChan(config)
}
func (b botAdapter) Send(c botapi.Chattable) (botapi.Message, error) {
	return b.Bot.Send(c)
}

func (b botAdapter) DeleteMessage(config botapi.DeleteMessageConfig) (botapi.APIResponse, error) {
	return b.Bot.DeleteMessage(config)
}
func (b botAdapter) UserName() string {
	return b.Bot.Self.UserName
}
