package config

type Config struct {
	BaseHeaders map[string]string
}

func NewConfig() *Config {
	return &Config{BaseHeaders: getHeaders()}
}

func getHeaders() map[string]string {
	return map[string]string{
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
}
