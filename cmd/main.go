package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	p "study/internal/parser"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	parentID := 1893

	// Загружаем или получаем куки
	cookieStr, err := p.loadCookiesFromFile()
	if err == nil && cookieStr != "" {
		currentCookieString = cookieStr
		// Пытаемся извлечь sessionToken, deviceID, userSessionID из файла кук
		// Для простоты можно переполучить через браузер, но пока оставим как есть
		// Лучше всё же один раз переполучить через браузер, чтобы гарантировать наличие этих полей
		fmt.Println("Куки загружены из файла cookies.txt")
		// Для надёжности всё равно обновим через браузер, чтобы извлечь значения
		fmt.Println("Обновляем куки через браузер для извлечения sessionToken, deviceID, userSessionID...")
		newCookieStr, err := fetchCookiesWithChrome()
		if err != nil {
			fmt.Printf("Предупреждение: не удалось обновить куки: %v\n", err)
		} else {
			currentCookieString = newCookieStr
		}
	} else {
		fmt.Println("Файл cookies.txt не найден или пуст. Получаем куки через браузер...")
		cookieStr, err = fetchCookiesWithChrome()
		if err != nil {
			fmt.Printf("Ошибка получения кук: %v\n", err)
			fmt.Println("Попробуйте вручную скопировать куки в файл cookies.txt или указать путь к браузеру в коде.")
			return
		}
		currentCookieString = cookieStr
	}

	categories, err := fetchCategoriesViaAPI(parentID)
	if err != nil {
		fmt.Printf("Ошибка при получении категорий: %v\n", err)
		return
	}
	fmt.Printf("Получено %d категорий\n", len(categories))

	allItems := make([]Item, 0)
	for _, cat := range categories {
		fmt.Printf("Сбор товаров для категории %d (%s)...\n", cat.ID, cat.Name)
		items, err := fetchAllItems(cat.ID)
		if err != nil {
			fmt.Printf("Ошибка для категории %d: %v\n", cat.ID, err)
			continue
		}
		fmt.Printf("Собрано %d товаров\n", len(items))
		allItems = append(allItems, items...)
		time.Sleep(1 * time.Second)
	}
	if err := saveItemsToFile(allItems, "items.json"); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	} else {
		fmt.Printf("Сохранено %d товаров в items.json\n", len(allItems))
	}
}
