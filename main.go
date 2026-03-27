package main

import (
	"log"
	"math/rand"
	"study/internal/api"
	"study/internal/bot"
	"study/internal/client"
	"study/internal/config"
	"study/internal/cookies"
	"study/internal/models"
	"study/internal/parser"
	"study/internal/storage"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config := config.NewConfig()

	var cookieProvider cookies.Provider
	fileProvider := cookies.NewFileProvider("cookies.txt")
	cookieStr, err := fileProvider.GetCookies()
	if err == nil && cookieStr != "" {
		cookieProvider = fileProvider
	} else {
		chromeProvider := cookies.NewChromeProvider()
		cookieProvider = chromeProvider
	}

	lentaClient := client.NewLentaClient(cookieProvider, config.BaseHeaders)
	categoryService := api.NewCategoryAPIService(lentaClient)
	itemService := api.NewItemAPIService(lentaClient)
	fileStorage := storage.NewFilesStorage("items.json")
	rootID := 1893

	parser := parser.NewParser(itemService, categoryService, fileStorage, rootID)
	parseFunc := func() ([]models.Item, error) {
		return parser.Run()
	}

	botInstance, err := bot.NewBot(parseFunc)
	if err != nil {
		log.Fatalf("Ошибка инит бота %v", err)
	}

	go func() {
		if err := botInstance.Start(); err != nil {
			log.Fatalf("Пизда копченому %v", err)
		}
	}()

	select {}
}
