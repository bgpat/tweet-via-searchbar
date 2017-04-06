package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bgpat/tweet-via-searchbar/middleware"
	"github.com/bgpat/tweet-via-searchbar/opensearch"
	"github.com/bgpat/twtr"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	consumer      = twtr.NewCredentials(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	redisURL      = os.Getenv("REDIS_URL")
	sessionSecret = os.Getenv("SESSION_SECRET")
	baseURL       = os.Getenv("BASE_URL")
)

type searchArgs struct {
	Query    string `form:"q" binding:"required"`
	Token    string `form:"token" binding:"required"`
	Secret   string `form:"secret" binding:"required"`
	Redirect string `form:"redirect" binding:"required"`
}

func main() {
	gob.Register(&twtr.Client{})
	r := gin.Default()
	store, err := sessions.NewRedisStore(100, "tcp", redisURL, "", []byte(sessionSecret))
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	r.Use(sessions.Sessions("session", store))
	r.Use(middleware.TwitterClient(consumer))

	r.LoadHTMLGlob("templates/*")

	r.Any("/", func(c *gin.Context) {
		client := middleware.Default(c)
		var user *twtr.User
		if client != nil && client.AccessToken != nil {
			user, _, err = client.VerifyCredentials(nil)
			if err != nil {
				printError(err)
			}
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"user":   user,
			"config": client.Config,
		})
	})

	r.GET("/opensearch.xml", func(c *gin.Context) {
		client := middleware.Default(c)
		if client.AccessToken == nil {
			c.XML(http.StatusUnauthorized, nil)
			return
		}
		user, _, err := client.VerifyCredentials(nil)
		if err != nil {
			printError(err)
		}
		o := opensearch.NewOpenSearch(user, client)
		c.Render(http.StatusOK, o)
	})

	r.POST("/search", func(c *gin.Context) {
		args := searchArgs{
			Redirect: "about:blank",
		}
		err := c.Bind(&args)
		if err != nil {
			printError(err)
		}
		client := twtr.NewClient(consumer, twtr.NewCredentials(args.Token, args.Secret))
		_, _, err = client.UpdateTweet(&twtr.Params{
			"status": args.Query,
		})
		if err != nil {
			printError(err)
		}
		c.Redirect(http.StatusFound, strings.Replace(args.Redirect, "{searchTerms}", args.Query, -1))
	})

	r.GET("/login", func(c *gin.Context) {
		client := middleware.Default(c)
		if client.AccessToken != nil {
			c.Redirect(http.StatusFound, baseURL)
			return
		}
		if client.RequestToken != nil {
			if verifier, ok := c.GetQuery("oauth_verifier"); ok {
				err := client.GetAccessToken(verifier)
				if err != nil {
					printError(err)
					c.Redirect(http.StatusFound, baseURL)
					return
				}
				client.Save()
				c.Redirect(http.StatusFound, baseURL)
				return
			}
			if denied, ok := c.GetQuery("dinied"); ok {
				durlStr := client.OAuthClient.ResourceOwnerAuthorizationURI + "?access_token=" + denied
				c.Redirect(http.StatusFound, durlStr)
				return
			}
		}
		urlStr, err := client.RequestTokenURL(baseURL + "/login")
		if err != nil {
			printError(err)
			c.Redirect(http.StatusFound, baseURL)
			return
		}
		client.Save()
		c.Redirect(http.StatusFound, urlStr)
	})

	r.GET("/logout", func(c *gin.Context) {
		s := sessions.Default(c)
		s.Clear()
		s.Save()
		c.Redirect(http.StatusFound, baseURL)
	})

	r.Run()
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v\n", err)
}
