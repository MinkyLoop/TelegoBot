package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"study/internal/cookies"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type LentaClient struct {
	httpClient     *http.Client
	cookieProvider cookies.Provider
	baseHeaders    map[string]string
	cookieString   string
}

func NewLentaClient(cookieProvider cookies.Provider, baseHeaders map[string]string) *LentaClient {
	client := &LentaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cookieProvider: cookieProvider,
		baseHeaders:    baseHeaders,
	}

	cookies, err := cookieProvider.GetCookies()
	if err == nil {
		client.cookieString = cookies
		client.updateHeadersFromCookies()
	}
	return client
}

func randomHex(length int) string {
	const charset = "0123456789abcdef"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (c *LentaClient) updateHeadersFromCookies() {
	if c.cookieString == "" {
		return
	}
	cookiesMap := parseCookies(c.cookieString)

	if val, ok := cookiesMap["sessiontoken"]; ok {
		c.baseHeaders["sessiontoken"] = val
	} else if val, ok := cookiesMap["Utk_SessionToken"]; ok {
		c.baseHeaders["sessiontoken"] = val
	}

	if val, ok := cookiesMap["deviceid"]; ok {
		c.baseHeaders["deviceid"] = val
		c.baseHeaders["x-device-id"] = val
	} else if val, ok := cookiesMap["Utk_DvcGuid"]; ok {
		c.baseHeaders["deviceid"] = val
		c.baseHeaders["x-device-id"] = val
	}

	if val, ok := cookiesMap["x-user-session-id"]; ok {
		c.baseHeaders["x-user-session-id"] = val
	} else if val, ok := cookiesMap["UserSessionId"]; ok {
		c.baseHeaders["x-user-session-id"] = val
	}
}

func parseCookies(cookieStr string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func (c *LentaClient) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.baseHeaders {
		req.Header.Set(k, v)
	}
	if c.cookieString != "" {
		req.Header.Set("Cookie", c.cookieString)
	}

	traceID := randomHex(32)
	spanID := randomHex(16)
	req.Header.Set("traceparent", fmt.Sprintf("00-%s-%s-01", traceID, spanID))
	req.Header.Set("x-trace-id", traceID)
	req.Header.Set("x-span-id", spanID)

	if req.Method == http.MethodGet {
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("priority", "u=0, i")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-site", "none")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Del("content-type")
	} else if req.Method == http.MethodPost {
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Del("x-requested-with")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Сессия устарела, обновляем куки...")
		newCookies, err := c.cookieProvider.GetCookies()
		if err != nil {
			return nil, fmt.Errorf("не удалось обновить куки: %w", err)
		}
		c.cookieString = newCookies
		c.updateHeadersFromCookies()

		newReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		if err != nil {
			return nil, err
		}
		newReq.Header = req.Header.Clone()
		newReq.Header.Set("Cookie", c.cookieString)

		for k, v := range c.baseHeaders {
			newReq.Header.Set(k, v)
		}
		return c.httpClient.Do(newReq)
	}
	return resp, nil
}

func (c *LentaClient) SetCookieString(cookie string) {
	c.cookieString = cookie
	c.updateHeadersFromCookies()
}
