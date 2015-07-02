package app

import (
	"github.com/go-martini/martini"
//"github.com/google/go-github/github"
	"github.com/cgcgbcbc/umeng2github/routes"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	goauth2 "golang.org/x/oauth2"
)

const (
	session_cookie = "secret"
	client = "5d88b7753ff5711224af"
	key = "85fae3cae0f1794d33a6a516bb0a75d4c46af2da"
)

func App() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte(session_cookie))))
	m.Use(oauth2.Github(
		&goauth2.Config{
			ClientID:     client,
			ClientSecret: key,
			Scopes:       []string{"user:email", "read:org"},
			RedirectURL:  "http://localhost:3000/oauth2callback",
		},
	))
	m.Get("/",oauth2.LoginRequired, routes.Index)
	m.Get("/dashboard", oauth2.LoginRequired, routes.Dashboard)
	return m
}