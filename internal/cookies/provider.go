package cookies

type Provider interface {
	GetCookies() (string, error)
}
