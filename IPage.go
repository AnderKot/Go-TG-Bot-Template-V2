package main

type IPage interface {
	// Common
	GetName() string

	// Input
	OnProcessingMessage(text string)
	OnProcessingKey(keyData string)

	// Navigation
	OnGetNextPage() IConstructor
	OnBackToParent() bool

	// Print
	GetMessageText() string
	GetKeyboard() IKeyboard
}
