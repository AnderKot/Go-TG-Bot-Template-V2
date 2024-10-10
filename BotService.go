package main

import (
	"Queue"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

func NewBotService(botAPI *tgBotAPI.BotAPI, database IRepository) BotService {
	return BotService{
		BotAPI:           botAPI,
		UserMutexes:      map[int64]*sync.Mutex{},
		UserPagesQueue:   map[int64]Queue.MyQueue[IPage]{},
		UserChatMessages: map[int64]Message{},
		Repository:       database,
	}
}

type Message struct {
	id       int
	text     string
	keyboard IKeyboard
}

type BotService struct {
	BotAPI           *tgBotAPI.BotAPI
	UserMutexes      map[int64]*sync.Mutex
	UserPagesQueue   map[int64]Queue.MyQueue[IPage]
	UserChatMessages map[int64]Message
	Repository       IRepository
}

func (bs *BotService) Start() {
	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 20

	//Получаем обновления от бота
	updates := bs.BotAPI.GetUpdatesChan(u)

	for u := range updates {
		// Обработка сообщения
		go func(bot *tgBotAPI.BotAPI, update tgBotAPI.Update) {
			chatID := GetChatIDFromUpdate(update)
			langCode := GetLangCodeFromUpdate(update)

			SaveLock(bs.UserMutexes[chatID])

			queue := bs.UserPagesQueue[chatID]
			exist, page := queue.Dequeue()
			if !exist {
				page = CreateMainMenu()
			}

			// Обработка сообщений
			if update.CallbackQuery != nil {
				(*page).OnProcessingKey(update.CallbackQuery.Data)
			}
			if update.Message != nil {
				(*page).OnProcessingMessage(update.Message.Text)
			}

			// Навигация по страницам
			nextPageConstructor := (*page).OnGetNextPage()
			if nextPageConstructor != nil {
				queue.Enqueue(page)
				p := nextPageConstructor.New()
				page = &(p)
			}

			if (*page).OnBackToParent() {
				oldExist, oldPage := queue.Dequeue()
				if oldExist {
					page = oldPage
				} else {
					page = CreateMainMenu()
				}
			}

			// Ответ пользователю
			newText := bs.Repository.GetTemplateText((*page).GetMessageText(), langCode)

			newKeyboard := (*page).GetKeyboard()

			oldMessage, isOldExist := bs.UserChatMessages[chatID]
			if isOldExist {
				if oldMessage.text != newText {
					isOk := bs.sendEdit(chatID, oldMessage.id, newText, bs.GenerateKeyboard(newKeyboard, langCode), false)
					if !isOk {
						panic("send edit message error")
					}
				}
			} else {
				isOk, messageID := bs.sendNew(chatID, newText, bs.GenerateKeyboard(newKeyboard, langCode), false)
				if isOk {
					bs.UserChatMessages[chatID] = Message{
						id:       messageID,
						text:     newText,
						keyboard: newKeyboard,
					}
				} else {
					panic("send new message error")
				}
			}

			queue.Enqueue(page)
			SaveUnlock(bs.UserMutexes[chatID])
		}(bs.BotAPI, u)
	}
}

func (bs *BotService) final() {

}

func (bs *BotService) sendNew(chat int64, text string, keyboard *tgBotAPI.InlineKeyboardMarkup, isMarkdown bool) (bool, int) {
	NewMsgRequest := tgBotAPI.NewMessage(chat, text)
	if isMarkdown {
		NewMsgRequest.ParseMode = "Markdown"
	}
	if keyboard != nil {
		if keyboard.InlineKeyboard != nil {
			NewMsgRequest.ReplyMarkup = keyboard
		}
	}
	NewMsgRespons, err := bs.BotAPI.Send(NewMsgRequest)
	if err != nil {
		return false, 0
	}

	return true, NewMsgRespons.MessageID
}

func (bs *BotService) sendEdit(chat int64, oldMessageID int, text string, keyboard *tgBotAPI.InlineKeyboardMarkup, isMarkdown bool) bool {
	NewMsgRequest := tgBotAPI.NewEditMessageText(chat, oldMessageID, text)
	if isMarkdown {
		NewMsgRequest.ParseMode = "Markdown"
	}
	if keyboard != nil {
		if keyboard.InlineKeyboard != nil {
			NewMsgRequest.ReplyMarkup = keyboard
		}
	}
	_, err := bs.BotAPI.Request(NewMsgRequest)
	if err != nil {
		return err.(*tgBotAPI.Error).Message == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message"
	}

	return true
}

func (bs *BotService) deleteMsg(chat int64, messageID int) {
	_, err := bs.BotAPI.Request(tgBotAPI.NewDeleteMessage(chat, messageID))
	if err != nil {
		return
	}
}

func (bs *BotService) GenerateKeyboard(kbd IKeyboard, langCode string) *tgBotAPI.InlineKeyboardMarkup {
	var botKeyboard tgBotAPI.InlineKeyboardMarkup
	for _, row := range kbd.GetRows() {
		var botRow []tgBotAPI.InlineKeyboardButton
		for _, key := range row.GetKeys() {
			botRow = append(botRow, tgBotAPI.NewInlineKeyboardButtonData(bs.Repository.GetTemplateText(key.GetTemplate(), langCode), key.GetData()))
		}
		if len(botRow) > 0 {
			botKeyboard.InlineKeyboard = append(botKeyboard.InlineKeyboard, botRow)
		}
	}
	if len(botKeyboard.InlineKeyboard) > 0 {
		return &botKeyboard
	}
	return nil
}
