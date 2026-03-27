package cookies

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ChromeProvider struct {
}

func NewChromeProvider() *ChromeProvider {
	return &ChromeProvider{}
}

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

func (c ChromeProvider) GetCookies() (string, error) {
	execPath := getBrowserPath()
	if execPath == "" {
		return "", fmt.Errorf("Нету браузера")
	}

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
		return "", err
	}

	var cookieParts []string
	for _, c := range cookiesData {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	cookieStr := strings.Join(cookieParts, "; ")

	if err = os.WriteFile("cookies.txt", []byte(cookieStr), 0644); err != nil {
		fmt.Printf("Предупреждение: не удалось сохранить куки в файл: %v\n", err)
	}
	fmt.Println("Куки успешно получены и сохранены в cookies.txt")

	return cookieStr, nil
}
