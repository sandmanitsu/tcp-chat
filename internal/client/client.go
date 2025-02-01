package client

type ClientConn chan string

type Client struct {
	Name     string
	CurrRoom string
}
