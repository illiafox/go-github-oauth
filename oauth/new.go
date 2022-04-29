package oauth

import (
	"fmt"

	"oauth/utils/config"
)

type Oauth struct {
	// AuthorizeURL redirects user to callback, ends with 'state=', so only random string is needed
	AuthorizeURL string

	// AuthorizeURL is used to get permanent user token,  ends with 'code=', so only code is needed
	AccessTokenURL string

	// GetUserURL is used to get user data (username, email, etc)
	// Pass token in Authorization
	GetUserURL string
}

func New(conf config.Oauth) *Oauth {

	return &Oauth{
		AuthorizeURL: fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&state=", conf.ClientID),

		AccessTokenURL: fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=",
			conf.ClientID, conf.ClientSecret),

		GetUserURL: "https://api.github.com/user",
	}

}
