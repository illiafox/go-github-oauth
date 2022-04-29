package oauth

import (
	"net/http"

	"go.uber.org/zap"
	"oauth/database"
)

func New(db *database.Database, logger *zap.Logger) http.Handler {
	root := http.NewServeMux()

	m := Methods{db, logger}

	root.HandleFunc("/callback", m.Callback)

	return root
}

type Methods struct {
	db     *database.Database
	logger *zap.Logger
}
