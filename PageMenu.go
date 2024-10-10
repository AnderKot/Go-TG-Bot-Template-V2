package main

import (
	"strconv"
)

// Констукторы страниц менюшек
func CreateMainMenu() *IPage {
	constructor := PageMenuConstructor{
		name:     "mainMenu",
		template: &PageTemplate{code: "mainMenuPage"},
		items: []MenuItem{
			{
				Name:        &MenuItemTemplate{code: "tournament"},
				Constructor: &tournamentPageMenuConstructor,
			}, {
				Name:        &MenuItemTemplate{code: "yourActivity"},
				Constructor: &PageMenuConstructor{},
			}, {
				Name:        &MenuItemTemplate{code: "profile"},
				Constructor: &profilePageMenuConstructor,
			},
		},
		isHasParent: false,
	}
	page := constructor.New()
	return &(page)
}

var tournamentPageMenuConstructor = PageMenuConstructor{
	name:     "tournamentMenu",
	template: &PageTemplate{code: "tournamentMenuPage"},
	items: []MenuItem{
		{
			Name:        &MenuItemTemplate{code: "create"},
			Constructor: &PageMenuConstructor{},
		}, {
			Name:        &MenuItemTemplate{code: "yours"},
			Constructor: &PageMenuConstructor{},
		}, {
			Name:        &MenuItemTemplate{code: "matchmaking"},
			Constructor: &PageMenuConstructor{},
		}, {
			Name:        &MenuItemTemplate{code: "running"},
			Constructor: &PageMenuConstructor{},
		}, {
			Name:        &MenuItemTemplate{code: "ended"},
			Constructor: &PageMenuConstructor{},
		},
	},
	isHasParent: true,
}

var profilePageMenuConstructor = PageMenuConstructor{
	name:        "profileMenu",
	template:    &PageTemplate{code: "profileMenuPage"},
	items:       []MenuItem{},
	isHasParent: true,
}

// Реализация страниц меню
type MenuItem struct {
	Name        ITemplate
	Constructor IConstructor
}

type PageMenuConstructor struct {
	name        string
	template    ITemplate
	items       []MenuItem
	isHasParent bool
}

func (lc *PageMenuConstructor) New() IPage {
	l := new(PageMenu)

	l.name = lc.name
	l.template = lc.template
	l.items = lc.items
	l.isHasParent = lc.isHasParent

	l.CreateKeyBoard()

	return l
}

type PageMenu struct {
	name           string
	template       ITemplate
	board          IKeyboard
	isHasParent    bool
	isNeedToParent bool

	items        []MenuItem
	selectedItem IConstructor
}

func (p *PageMenu) CreateKeyBoard() {
	kb := Keyboard{Rows: make([]IKeyRow, 0)}
	kbr := KeyRow{make([]IKey, 0)}

	for index, item := range p.items {
		kbr.Keys = append(kbr.Keys, &Key{Name: item.Name, Data: strconv.Itoa(index)})
		if (index+1)%2 == 0 {
			kb.Rows = append(kb.Rows, &kbr)
			kbr = KeyRow{make([]IKey, 0)}
		}
	}
	kb.Rows = append(kb.Rows, &kbr)

	if p.isHasParent {
		kb.Rows = append(kb.Rows, &KeyRow{Keys: []IKey{
			Key{Name: &onBackToParentTemplate{}, Data: onBackToParent},
		}})
	}

	p.board = &kb
}

type PageTemplate struct {
	code string
}

func (pt *PageTemplate) isTranslated() bool { return true }

func (pt *PageTemplate) GetTemplateText() string { return "" }

func (pt *PageTemplate) GetTemplateCode() string { return pt.code }

type MenuItemTemplate struct {
	code string
}

func (pt *MenuItemTemplate) isTranslated() bool { return true }

func (pt *MenuItemTemplate) GetTemplateText() string { return "" }

func (pt *MenuItemTemplate) GetTemplateCode() string { return pt.code }

type onBackToParentTemplate struct{}

func (obpt *onBackToParentTemplate) isTranslated() bool { return true }

func (obpt *onBackToParentTemplate) GetTemplateText() string { return "onBackToParent" }

func (obpt *onBackToParentTemplate) GetTemplateCode() string { return onBackToParent }

// IPage >>
// Common
func (p *PageMenu) GetName() string {
	return p.name
}

// Input
func (p *PageMenu) OnProcessingMessage(text string) {
}

func (p *PageMenu) OnProcessingKey(keyData string) {
	switch keyData {
	case onBackToParent:
		{
			p.isNeedToParent = true
		}
	default:
		{
			index, _ := strconv.Atoi(keyData)
			p.selectedItem = p.items[index].Constructor
		}
	}
}

// Navigation
func (p *PageMenu) OnGetNextPage() IConstructor {
	return p.selectedItem
}

func (p *PageMenu) OnBackToParent() bool {
	return p.isNeedToParent
}

// Print
func (p *PageMenu) GetMessageText() ITemplate {
	return p.template
}

func (p *PageMenu) GetKeyboard() IKeyboard {
	return p.board
}

// IPage <<
