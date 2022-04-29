package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"oauth/database"
	"oauth/server/handlers/oauth"
	"oauth/server/handlers/site"
	"oauth/utils/config"
)

func New(logger *zap.Logger, db *database.Database, conf config.Host) *http.Server {
	root := http.NewServeMux()

	root.Handle("/metrics", promhttp.Handler())

	root.Handle("/oauth/", http.StripPrefix("/oauth", oauth.New(db, logger)))

	root.Handle("/", site.New(db, logger))

	//

	return &http.Server{
		Addr: "0.0.0.0:" + conf.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      root,
	}
}
