package main

type IKeyboard interface {
	GetRows() []IKeyRow
}

type Keyboard struct {
	Rows []IKeyRow
}

func (k *Keyboard) GetRows() []IKeyRow {
	return k.Rows
}

type IKeyRow interface {
	GetKeys() []IKey
}

type KeyRow struct {
	Keys []IKey
}

func (k *KeyRow) GetKeys() []IKey {
	return k.Keys
}

type IKey interface {
	GetTemplate() ITemplate
	GetData() string
}

type Key struct {
	Name ITemplate
	Data string
}

func (k Key) GetTemplate() ITemplate {
	return k.Name
}

func (k Key) GetData() string {
	return k.Data
}
