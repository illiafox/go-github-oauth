package site

import (
	"net/http"

	"go.uber.org/zap"
	"oauth/oauth"
	"oauth/server/cookie"
	"oauth/utils/templates"
)

func (m Methods) toLogin(w http.ResponseWriter, r *http.Request) {
	state := oauth.Generate()

	err := m.db.Memcached.StoreState(state)
	if err != nil {
		m.logger.Error("oauth: memcached: store state",
			zap.String("state", state),
			zap.Error(err),
		)
		templates.Message.Internal(w)

		return
	}

	http.Redirect(w, r, m.db.Oauth.AuthorizeURL+state, http.StatusTemporaryRedirect)

}

func (m Methods) Index(w http.ResponseWriter, r *http.Request) {

	token, err := cookie.Get(r, w)

	if err != nil || token == "" {

		m.toLogin(w, r)

		return
	}

	id, err := m.db.Postgres.Session.Exists(token)

	if err != nil {
		m.logger.Error("index: session exists",
			zap.String("token", token),
			zap.Error(err),
		)
		templates.Message.Internal(w)

		return
	}

	if id < 0 {
		m.toLogin(w, r)

		return
	}

	username, err := m.db.Postgres.User.Username(id)
	if err != nil {
		m.logger.Error("index: get username", zap.Int64("user_id", id), zap.Error(err))
		templates.Message.Internal(w)

		return
	}

	templates.Main.Execute(w, username)
}
