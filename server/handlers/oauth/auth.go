package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"oauth/oauth"
	"oauth/server/cookie"
	"oauth/utils/templates"
)

func (m Methods) Callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	state := query.Get("state")
	if state == "" {
		templates.Message.Execute(w, "state not found")

		return
	}

	// get code
	code := query.Get("code")
	if state == "" {
		templates.Message.Execute(w, "code not found")

		return
	}

	// find state in cache
	exist, err := m.db.Memcached.LookupState(state)
	if err != nil { // only internal
		m.logger.Error("oauth: memcached: lookup state",
			zap.Error(err),
			zap.String("state", state),
		)
		templates.Message.Internal(w)

		return
	}

	if !exist {
		templates.Message.Execute(w, "state not found")

		return
	}

	// // //

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get token
	request, err := http.NewRequestWithContext(ctx, "GET", m.db.Oauth.AccessTokenURL+code, nil)

	if err != nil {
		m.logger.Error("oauth: NewRequestWithContext", zap.Error(err), zap.String("url", m.db.Oauth.AccessTokenURL+code))
		templates.Message.Internal(w)

		return
	}

	// Do request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		m.logger.Error("oauth: do request", zap.Error(err))
		templates.Message.Internal(w)

		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		m.logger.Error("oauth: reading body", zap.Error(err))
		templates.Message.Internal(w)

		return
	}
	r.Body.Close()

	values, err := url.ParseQuery(string(data))
	if err != nil {
		m.logger.Error("oauth: parsing query",
			zap.Error(err),
			zap.String("query", string(data)),
		)
		templates.Message.Internal(w)

		return
	}

	token := values.Get("access_token")
	if token == "" {
		templates.Message.Execute(w, "token error<br>"+values.Get("error_description"))

		return
	}

	// // //

	// get username
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err = http.NewRequestWithContext(ctx, "GET", m.db.Oauth.GetUserURL, nil)

	// set token as authorization
	request.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		m.logger.Error("oauth: NewRequestWithContext",
			zap.Error(err),
			zap.String("url", m.db.Oauth.GetUserURL),
		)
		templates.Message.Internal(w)

		return
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		m.logger.Error("oauth: do request", zap.Error(err))
		templates.Message.Internal(w)

		return
	}

	var user getToken

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		m.logger.Error("oauth: decode json", zap.Error(err))
		templates.Message.Internal(w)

		return
	}

	if user.Message != "" {
		m.logger.Error("oauth: get username", zap.String("error", user.Message))
		templates.Message.Internal(w)

		return
	}

	if user.Login == "" {
		templates.Message.Execute(w, "couldn't get github username")

		return
	}

	username, err := m.db.Postgres.User.Exists(user.ID)
	if err != nil {
		m.logger.Error("oauth: user exists",
			zap.Error(err),
			zap.Int64("user_id", user.ID),
		)
		templates.Message.Internal(w)

		return
	}

	if username == "" {
		err = m.db.Postgres.User.Create(user.ID, token, user.Login)
		if err != nil {
			m.logger.Error("oauth: create user",
				zap.Int64("user_id", user.ID),
				zap.String("username", user.Login),
				zap.String("token", token),
				zap.Error(err),
			)
			templates.Message.Internal(w)

			return
		}
	}

	if username != user.Login {
		err = m.db.Postgres.User.UpdateUsername(user.ID, user.Login)
		if err != nil {
			m.logger.Error("oauth: update username",
				zap.Int64("user_id", user.ID),
				zap.String("username", user.Login),
				zap.Error(err),
			)
			templates.Message.Internal(w)

			return
		}
	}

	// // //

	// new token for session
	token = oauth.Generate()

	err = cookie.Set(r, w, token)
	if err != nil {
		templates.Message.Execute(w, "cookie error!<br>please,enable them")

		return
	}

	err = m.db.Postgres.Session.New(token, user.ID)
	if err != nil {
		m.logger.Error("oauth: new session",
			zap.Int64("user_id", user.ID),
			zap.String("token", token),
			zap.Error(err),
		)
		templates.Message.Internal(w)

		return
	}

	templates.Message.Execute(w, `Redirecting to main page
    <meta http-equiv="refresh" content="2 url=/">`)
}

type getToken struct {
	Login   string `json:"login"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}
