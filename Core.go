package main

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

// Const Keys Data
const (
	onBackToParent string = "onBackToParent"
)

type NewPageConstructor func(_arg any) IPage

func GetChatIDFromUpdate(update tgBotAPI.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func GetLangCodeFromUpdate(update tgBotAPI.Update) string {
	if update.Message != nil {
		return update.Message.From.LanguageCode
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LanguageCode
	}

	return "en"
}

func SaveLock(m *sync.Mutex) {
	if m == nil {
		m = &sync.Mutex{}
	}
	m.Lock()
}

func SaveUnlock(m *sync.Mutex) {
	if m != nil {
		m.Unlock()
	}
}
