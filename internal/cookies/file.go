package cookies

import (
	"os"
	"strings"
)

type FileProvider struct {
	name string
}

func NewFileProvider(name string) *FileProvider {
	return &FileProvider{name}
}
func (p *FileProvider) GetCookies() (string, error) {
	data, err := os.ReadFile(p.name)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
