package cookie

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	cookie = "4DB2E134E453D713A48316A422DBDB812F4C79C2815F152C0147A3CB86864D53"
	name   = "credentials" // cookie name
	key    = "session"     // map key with token
)

var (
	store *sessions.CookieStore
)

func init() {
	store = sessions.NewCookieStore([]byte(cookie))
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode
	store.Options.MaxAge = 60 * 60 * 24 * 28 // 28 days
}

func Set(r *http.Request, w http.ResponseWriter, token string) error {
	session, err := store.Get(r, name)
	if err != nil {
		return err
	}

	session.Values[key] = token

	return sessions.Save(r, w)
}

func Get(r *http.Request, w http.ResponseWriter) (string, error) {
	session, err := store.Get(r, name)
	if err != nil {
		return "", err
	}

	token, ok := session.Values[key].(string)
	if !ok || token == "" {
		return "", nil
	}

	return token, nil
}

func Delete(r *http.Request, w http.ResponseWriter) error {
	session, err := store.Get(r, name)
	if err != nil {
		return err
	}

	delete(session.Values, key)

	return sessions.Save(r, w)
}
