package parser

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

	s "study/internal/parser/struct"
)

var baseHeaders = map[string]string{
	"accept":              "application/json",
	"accept-language":     "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
	"client":              "angular_web_0.0.2",
	"content-type":        "application/json",
	"origin":              "https://lenta.com",
	"sec-ch-ua":           `"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`,
	"sec-ch-ua-mobile":    "?0",
	"sec-ch-ua-platform":  `"Windows"`,
	"user-agent":          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
	"x-delivery-mode":     "pickup",
	"x-device-os":         "Web",
	"x-device-os-version": "12.4.8",
	"x-domain":            "moscow",
	"x-platform":          "omniweb",
	"x-retail-brand":      "lo",
}

var (
	currentCookieString string
	sessionToken        string
	deviceID            string
	userSessionID       string
)

func getBrowserPath() string {
	chromePaths := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
	}
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		chromePaths = append(chromePaths, filepath.Join(localAppData, `Google\Chrome\Application\chrome.exe`))
	}
	for _, p := range chromePaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func fetchCookiesWithChrome() (string, error) {
	execPath := getBrowserPath()
	if execPath == "" {
		return "", fmt.Errorf("не найден поддерживаемый браузер (Chrome или Яндекс). Установите путь вручную в коде.")
	}
	fmt.Printf("Используем браузер: %s\n", execPath)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.ExecPath(execPath),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("headless", true),
	)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var cookiesData []*network.Cookie

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://lenta.com"),
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookiesData, err = network.GetCookies().Do(ctx)

			return err
		}),
	)
	if err != nil {
		return "", fmt.Errorf("chromedp ошибка: %w", err)
	}

	var cookieParts []string
	for _, c := range cookiesData {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
		switch c.Name {
		case "sessiontoken", "Utk_SessionToken":
			sessionToken = c.Value
		case "deviceid", "Utk_DvcGuid":
			deviceID = c.Value
		case "x-user-session-id", "UserSessionId":
			userSessionID = c.Value
		}
	}
	cookieStr := strings.Join(cookieParts, "; ")

	if sessionToken == "" {
		for _, c := range cookiesData {
			if strings.Contains(strings.ToLower(c.Name), "session") {
				sessionToken = c.Value
				break
			}
		}
	}
	if deviceID == "" {
		for _, c := range cookiesData {
			if strings.Contains(strings.ToLower(c.Name), "device") {
				deviceID = c.Value
				break
			}
		}
	}
	if userSessionID == "" {
		for _, c := range cookiesData {
			if strings.Contains(strings.ToLower(c.Name), "user") && strings.Contains(strings.ToLower(c.Name), "session") {
				userSessionID = c.Value
				break
			}
		}
	}

	err = os.WriteFile("cookies.txt", []byte(cookieStr), 0644)
	if err != nil {
		fmt.Printf("Предупреждение: не удалось сохранить куки в файл: %v\n", err)
	}
	fmt.Println("Куки успешно получены и сохранены в cookies.txt")
	return cookieStr, nil
}

func loadCookiesFromFile() (string, error) {
	data, err := os.ReadFile("cookies.txt")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func randomHex(length int) string {
	const charset = "0123456789abcdef"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func buildRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}

	if currentCookieString != "" {
		req.Header.Set("Cookie", currentCookieString)
	}

	if sessionToken != "" {
		req.Header.Set("sessiontoken", sessionToken)
	}
	if deviceID != "" {
		req.Header.Set("deviceid", deviceID)
		req.Header.Set("x-device-id", deviceID)
	}
	if userSessionID != "" {
		req.Header.Set("x-user-session-id", userSessionID)
	}

	traceID := randomHex(32)
	spanID := randomHex(16)
	req.Header.Set("traceparent", fmt.Sprintf("00-%s-%s-01", traceID, spanID))
	req.Header.Set("x-trace-id", traceID)
	req.Header.Set("x-span-id", spanID)

	if method == http.MethodPost {
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
	} else if method == http.MethodGet {
		req.Header.Del("content-type")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("priority", "u=0, i")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-site", "none")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("upgrade-insecure-requests", "1")
	}

	return req, nil
}

func fetchCategoriesViaAPI(parentID int) ([]s.Category, error) {
	url := fmt.Sprintf("https://lenta.com/api-gateway/v1/catalog/categories?parentId=%d&depth=2", parentID)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := buildRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Сессия устарела, обновляем куки...")
		newCookies, err := fetchCookiesWithChrome()
		if err != nil {
			return nil, fmt.Errorf("не удалось обновить куки: %w", err)
		}
		currentCookieString = newCookies
		return fetchCategoriesViaAPI(parentID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("не удалось получить категории: %s, тело: %s", resp.Status, string(body))
	}
	var categoriesResp s.CategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&categoriesResp); err != nil {
		return nil, err
	}
	return categoriesResp.Categories, nil
}

func fetchAllItems(categoryID int) ([]s.Item, error) {
	const limit = 40
	offset := 0
	var allItems []s.Item
	client := &http.Client{Timeout: 30 * time.Second}
	url := "https://lenta.com/api-gateway/v1/catalog/items"

	for {
		reqBody := s.ItemsRequest{
			CategoryID: categoryID,
			Filters:    s.Filters{Checkbox: []interface{}{}, Multicheckbox: []interface{}{}, Range: []interface{}{}},
			Sort:       s.Sort{Type: "popular", Order: "desc"},
			Limit:      limit,
			Offset:     offset,
		}
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := buildRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Сессия устарела, обновляем куки...")
			newCookies, err := fetchCookiesWithChrome()
			if err != nil {
				return nil, fmt.Errorf("не удалось обновить куки: %w", err)
			}
			currentCookieString = newCookies
			return fetchAllItems(categoryID)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("статус %s, тело: %s", resp.Status, string(body))
		}

		var respData s.ItemsResponse
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return nil, err
		}

		for i := range respData.Items {
			respData.Items[i].CategoryID = categoryID
		}
		allItems = append(allItems, respData.Items...)

		if len(respData.Items) < limit {
			break
		}
		offset += limit
		time.Sleep(500 * time.Millisecond)
	}
	return allItems, nil
}

func saveItemsToFile(items []s.Item, filename string) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
