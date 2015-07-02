package routes

import (
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func Dashboard(tokens oauth2.Tokens, r render.Render, s sessions.Session) {
	if s.Get("key") == nil {
		r.Redirect("/", 302)
	}
	data := map[string]interface{}{"Title": "Umeng to Github"}
	r.HTML(200, "dashboard", data)
}

func Index(s sessions.Session, r render.Render) {
	s.Set("key", "value")
	r.Redirect("/dashboard", 302)
}
