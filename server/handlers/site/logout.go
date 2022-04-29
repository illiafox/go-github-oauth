package site

import (
	"net/http"

	"go.uber.org/zap"
	"oauth/server/cookie"
	"oauth/utils/templates"
)

func (m Methods) Logout(w http.ResponseWriter, r *http.Request) {

	token, err := cookie.Get(r, w)

	if err != nil || token == "" {

		m.toLogin(w, r)

		return
	}

	id, err := m.db.Postgres.Session.Exists(token)

	if err != nil {
		m.logger.Error("logout: session exists",
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

	err = m.db.Postgres.Session.Delete(token)
	if err != nil {
		m.logger.Error("logout: delete session",
			zap.String("token", token),
			zap.Error(err),
		)
		templates.Message.Internal(w)

		return
	}

	templates.Message.Execute(w, "Logged out")
}
