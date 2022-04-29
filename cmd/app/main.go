package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"oauth/database"
	"oauth/server"
	"oauth/utils/config"
	zaplog "oauth/utils/zap"
)

func main() {

	var (
		conf      *config.Config
		logger    *zap.Logger
		forceHTTP = flag.Bool("http", false, "force forceHTTP mode")
	)

	{
		logfile := flag.String("log", "log.txt", "log file path")
		path := flag.String("config", "config.toml", "config file path")
		env := flag.Bool("env", false, "read from environment")
		flag.Parse()

		file, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE, 0655)
		if err != nil {
			log.Fatalln("opening log file: ", err)
		}

		logger = zaplog.NewLogger(file)

		conf, err = config.Read(*path, *env)
		if err != nil {
			logger.Fatal("reading config", zap.Error(err))
		}
	}

	logger.Info("Initializing database")

	db, err := database.New(conf)
	if err != nil {
		logger.Fatal("initializing:", zap.Error(err))
	}

	defer db.Close(logger)

	logger.Info("Clearing old sessions")
	{
		deleted, err := db.Postgres.Session.DeleteOld(30)
		if err != nil {
			logger.Error("clearing sessions:", zap.Error(err))

			return
		}

		logger.Info("Done", zap.Int64("deleted", deleted))
	}

	srv := server.New(logger, db, conf.Host)

	ch := make(chan os.Signal, 1)

	go func() {
		logger.Info("Server started at " + srv.Addr)

		if *forceHTTP {
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS(conf.Host.Cert, conf.Host.Key)
		}

		if err != nil {
			if err == http.ErrServerClosed {
				logger.Info("Done")
			} else {
				logger.Error("Server:", zap.Error(err))
			}
			ch <- nil
		}
	}()

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)

	<-ch

	// Create a deadline to wait for closing all connections
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	logger.Info("Shutting down server")
	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Shutting:", zap.Error(err))
	}

}
