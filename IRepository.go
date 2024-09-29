package main

import (
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type IRepository interface {
	GetTemplateText(name string, langCode string) string
	ExportTemplates()
}

type User struct {
	ID           int64 `gorm:"primaryKey;autoIncrement:false"`
	NickName     string
	LanguageCode string `gorm:"default:'en'"`
}

type Templates struct {
	Name string `gorm:"primaryKey;" csv:"name"`
	Code string `gorm:"primaryKey;" csv:"code"`
	Text string `csv:"text"`
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
	file, err := os.Open(filepath.Join("Core", "templates.csv"))
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

func (r *Repository) GetTemplateText(name string, code string) string {
	templates := new(Templates)
	templates.Name = name
	templates.Code = code
	r.Database.Where(&templates).Take(templates)
	if templates.Text != "" {
		return templates.Text
	}
	templates.Code = "en"
	r.Database.Where(&templates).Take(templates)
	if templates.Text == "" {
		templates.Text = "Template name not specificate"
		r.Database.Save(templates)
		templates.Code = code
		r.Database.Save(templates)
	}
	return templates.Text
}

func (r *Repository) ExportTemplates() {

}
