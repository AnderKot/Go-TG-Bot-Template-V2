package main

import (
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type IRepository interface {
	GetUserByChatID(id int64) User
	GetTemplateText(template ITemplate, langCode string) string
	ExportTemplates()
}

type User struct {
	ID           int64 `gorm:"primaryKey;autoIncrement:false"`
	NickName     string
	LanguageCode string `gorm:"default:'en'"`
}

type Templates struct {
	Code     string `gorm:"primaryKey;" csv:"code"`
	langCode string `gorm:"primaryKey;" csv:"lang_code"`
	Text     string `csv:"text"`
}

func InitBase(dialector gorm.Dialector) IRepository {

	db, err := gorm.Open(dialector, &gorm.Config{
		/*
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "core.",
				SingularTable: false,
			}
		*/
	})
	if err != nil {
		panic("failed to connect database")
	}

	//db.Exec("CREATE SCHEMA IF NOT EXISTS core AUTHORIZATION " + dbLogin + ";")

	//db.AutoMigrate(&User{})
	db.AutoMigrate(&Templates{})
	//db.AutoMigrate(&Message{})

	// template migration
	file, err := os.Open(filepath.Join("", "templates.csv"))
	if err == nil {
		defer file.Close()
		var templates []Templates
		err = gocsv.Unmarshal(file, &templates)
		if err != nil {
			panic(err)
		}
		db.Create(templates)
	}

	//db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Templates{})
	return &(Repository{Database: db})
}

type Repository struct {
	Database *gorm.DB
}

func (r *Repository) GetUserByChatID(id int64) User {
	user := new(User)
	user.ID = id
	r.Database.Where(&user).Take(user)
	return *user
}

func (r *Repository) GetTemplateText(template ITemplate, langCode string) string {
	if template.isTranslated() {
		templates := new(Templates)
		templates.Code = template.GetTemplateCode()
		templates.langCode = langCode
		r.Database.Where(&templates).Take(templates)
		if templates.Text != "" {
			return templates.Text
		}
		templates.langCode = "en"
		r.Database.Where(&templates).Take(templates)
		if templates.Text == "" {
			templates.Text = "Template name not specificate"
			r.Database.Save(templates)
			templates.langCode = langCode
			r.Database.Save(templates)
		}
		return templates.Text
	} else {
		return template.GetTemplateText()
	}
}

func (r *Repository) ExportTemplates() {

}
