package bot

import (
	"log"
	"sub-cos-counter/internal/config"
	"sub-cos-counter/internal/services"
	"time"

	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot                 *telebot.Bot
	subscriptionService *services.SubscriptionService
	analyticsService    *services.AnalyticsService
	userStates          map[int64]*UserState
}

type UserState struct {
	State string
	Data  map[string]interface{}
}

const (
	StateIdle                = "idle"
	StateAddingSubscription  = "adding_subscription"
	StateWaitingForName      = "waiting_for_name"
	StateWaitingForCost      = "waiting_for_cost"
	StateWaitingForDate      = "waiting_for_date"
)

func NewBot(cfg *config.Config, subscriptionService *services.SubscriptionService, analyticsService *services.AnalyticsService) (*Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.GetBotToken(),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		bot:                 bot,
		subscriptionService: subscriptionService,
		analyticsService:    analyticsService,
		userStates:          make(map[int64]*UserState),
	}

	b.setupHandlers()
	return b, nil
}

func (b *Bot) Start() {
	log.Println("Bot started...")
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) getUserState(userID int64) *UserState {
	if state, exists := b.userStates[userID]; exists {
		return state
	}
	
	state := &UserState{
		State: StateIdle,
		Data:  make(map[string]interface{}),
	}
	b.userStates[userID] = state
	return state
}

func (b *Bot) setState(userID int64, state string) {
	userState := b.getUserState(userID)
	userState.State = state
}

func (b *Bot) setData(userID int64, key string, value interface{}) {
	userState := b.getUserState(userID)
	userState.Data[key] = value
}

func (b *Bot) getData(userID int64, key string) interface{} {
	userState := b.getUserState(userID)
	return userState.Data[key]
}

func (b *Bot) clearUserState(userID int64) {
	b.userStates[userID] = &UserState{
		State: StateIdle,
		Data:  make(map[string]interface{}),
	}
}