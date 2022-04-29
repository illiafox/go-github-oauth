package site

import (
	"net/http"

	"go.uber.org/zap"
	"oauth/database"
)

func New(db *database.Database, logger *zap.Logger) http.Handler {
	root := http.NewServeMux()

	m := Methods{db, logger}

	root.HandleFunc("/logout", m.Logout)
	root.HandleFunc("/", m.Index)

	return root
}

type Methods struct {
	db     *database.Database
	logger *zap.Logger
}
