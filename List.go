package main

import (
	"strconv"
	"strings"
)

const (
	onNextListPage string = "onNextListPage"
	onPrefListPage string = "onPrefListPage"
)

type ListItem struct {
	Name        string
	Description string
	Constructor IConstructor
}

type Lister interface {
	GetListItems(no int) []ListItem
	GetNextIsExists(no int) bool
}

type ListConstructor struct {
	name                string
	repository          IRepository
	nextPageConstructor IConstructor
	lister              Lister
	columns             int
}

func (lc *ListConstructor) New() IPage {
	l := new(List)

	l.name = lc.name
	l.repository = lc.repository
	l.lister = lc.lister
	l.columns = lc.columns

	l.CreateKeyBoardAndMessage()

	return l
}

type List struct {
	name           string
	repository     IRepository
	messageBuilder strings.Builder
	board          IKeyboard
	isNeedToParent bool

	lister       Lister
	items        []ListItem
	selectedItem IConstructor
	curItemsNo   int

	columns int
}

func (l *List) CreateKeyBoardAndMessage() {
	l.messageBuilder.Reset()
	//l.messageText = l.repository.GetTemplate(l.name)
	kb := Keyboard{Rows: make([]IKeyRow, 0)}
	kbr := KeyRow{make([]IKey, 0)}

	if l.lister.GetNextIsExists(l.curItemsNo) {
		l.items = l.lister.GetListItems(l.curItemsNo)
		for index, item := range l.items {
			l.items = append(l.items)
			strIndex := "[" + strconv.Itoa(index+1) + "]"
			_, _ = l.messageBuilder.WriteString("\n" + strIndex + ": " + item.Description)
			kbr.Keys = append(kbr.Keys, &Key{Name: ListItemKeyTemplate{text: strIndex}, Data: item.Name})
			if (index+1)%l.columns == 0 {
				kb.Rows = append(kb.Rows, &kbr)
				kbr = KeyRow{make([]IKey, 0)}
				_, _ = l.messageBuilder.WriteString("\n")
			}
		}
		kbr = KeyRow{
			Keys: []IKey{
				Key{
					Name: onNextListPageTemplate{},
					Data: onNextListPage,
				}}}
		kb.Rows = append(kb.Rows, &kbr)
	}

	if l.curItemsNo > 0 {
		kbr.Keys = append(kbr.Keys, Key{
			Name: onPrefListPageTemplate{},
			Data: onPrefListPage,
		})
	}

	kb.Rows = append(kb.Rows, &KeyRow{Keys: []IKey{
		Key{Name: &onBackToParentTemplate{}, Data: onBackToParent},
	}})

	l.board = &kb
}

type ListItemKeyTemplate struct {
	text string
}

func (obpt ListItemKeyTemplate) isTranslated() bool { return false }

func (obpt ListItemKeyTemplate) GetTemplateText() string { return obpt.text }

func (obpt ListItemKeyTemplate) GetTemplateCode() string { return "" }

type onNextListPageTemplate struct{}

func (obpt onNextListPageTemplate) isTranslated() bool { return true }

func (obpt onNextListPageTemplate) GetTemplateText() string { return "" }

func (obpt onNextListPageTemplate) GetTemplateCode() string { return onNextListPage }

type onPrefListPageTemplate struct{}

func (obpt onPrefListPageTemplate) isTranslated() bool { return true }

func (obpt onPrefListPageTemplate) GetTemplateText() string { return "" }

func (obpt onPrefListPageTemplate) GetTemplateCode() string { return onPrefListPage }

// IPage >>
// Common
func (l *List) GetName() string {
	return l.name
}

// Input
func (l *List) OnProcessingMessage(text string) {
}

func (l *List) OnProcessingKey(keyData string) {
	switch keyData {
	case onBackToParent:
		{
			l.isNeedToParent = true
		}
	case onNextListPage:
		{
			l.curItemsNo++
			l.CreateKeyBoardAndMessage()
		}
	case onPrefListPage:
		{
			if l.curItemsNo > 0 {
				l.curItemsNo--
				l.CreateKeyBoardAndMessage()
			}
		}
	default:
		{
			index, _ := strconv.Atoi(keyData)
			l.selectedItem = l.items[index].Constructor
		}
	}
}

// Navigation
func (l *List) OnGetNextPage() IConstructor {
	return l.selectedItem
}

func (l *List) OnBackToParent() bool {
	return l.isNeedToParent
}

// Print
func (l *List) GetMessageText() ITemplate {
	return &ListTemplate{text: l.messageBuilder.String()}
}

func (l *List) GetKeyboard() IKeyboard {
	return l.board
}

// IPage <<

type ListTemplate struct {
	text string
}

func (obpt *ListTemplate) isTranslated() bool { return false }

func (obpt *ListTemplate) GetTemplateText() string { return obpt.text }

func (obpt *ListTemplate) GetTemplateCode() string { return "" }
