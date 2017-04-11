package middleware

import (
	"encoding/gob"

	"github.com/bgpat/twtr"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	DefaultKey = "github.com/bgpat/twtr/gin/middleware"
)

type Client struct {
	*twtr.Client
	Context *gin.Context
	Config  Config
}

type Config struct {
	Redirect string `form:"redirect" binding:"omitempty"`
}

func init() {
	gob.Register(&Client{})
}

func NewClient(c *twtr.Client, ctx *gin.Context) *Client {
	return &Client{
		Client:  c,
		Context: ctx,
		Config: Config{
			Redirect: "about:blank",
		},
	}
}

func (c *Client) Save() {
	s := sessions.Default(c.Context)
	ctx := c.Context
	c.Context = nil
	s.Set(DefaultKey, c)
	s.Save()
	c.Context = ctx
}

func TwitterClient(consumer *twtr.Credentials) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		v := s.Get(DefaultKey)
		var client *Client
		if v == nil {
			client = NewClient(twtr.NewClient(consumer, nil), c)
			client.Save()
		} else {
			client = v.(*Client)
			client.Context = c
		}
		c.Set(DefaultKey, client)
		c.Next()
	}
}

func Default(c *gin.Context) *Client {
	return c.MustGet(DefaultKey).(*Client)
}
