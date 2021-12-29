package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/honoroit/cache"
	"gitlab.com/etke.cc/honoroit/config"
	"gitlab.com/etke.cc/honoroit/logger"
	"gitlab.com/etke.cc/honoroit/matrix"
)

const (
	enableEncryption = false
	fatalmessage     = "recovery(): %v"
)

var (
	version = "development"
	bot     *matrix.Bot
	log     *logger.Logger
)

func main() {
	var err error
	cfg := config.New()
	log = logger.New("honoroit.", cfg.LogLevel)
	defer recovery(cfg.RoomID)

	log.Info("#############################")
	log.Info("Honoroit " + version)
	log.Info("Matrix: true")
	log.Info("#############################")

	initBot(cfg)
	initShutdown()

	log.Debug("starting bot...")
	if err = bot.Start(); err != nil {
		// nolint // Fatal = panic, not os.Exit()
		log.Fatal("matrix bot crashed: %v", err)
	}
}

func initBot(cfg *config.Config) {
	db, err := sql.Open(cfg.DB.Dialect, cfg.DB.DSN)
	if err != nil {
		log.Fatal("cannot initialize SQL database: %v", err)
	}
	inmemoryCache := cache.New(time.Duration(cfg.TTL) * time.Minute)
	botConfig := &matrix.Config{
		Homeserver: cfg.Homeserver,
		Login:      cfg.Login,
		Password:   cfg.Password,
		Token:      cfg.Token,
		LogLevel:   cfg.LogLevel,
		RoomID:     cfg.RoomID,
		Text:       (*matrix.Text)(&cfg.Text),
		Cache:      inmemoryCache,
		DB:         db,
		Dialect:    cfg.DB.Dialect,
	}
	bot, err = matrix.NewBot(botConfig)
	if err != nil {
		// nolint // Fatal = panic, not os.Exit()
		log.Fatal("cannot create the matrix bot: %v", err)
	}
	log.Debug("bot has been created")

	if enableEncryption {
		if err = bot.WithEncryption(); err != nil {
			// nolint // Fatal = panic, not os.Exit()
			log.Fatal("cannot initialize e2ee support: %v", err)
		}
		log.Debug("end-to-end encryption support initialized")
	}
}

func initShutdown() {
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for range listener {
			bot.Stop()
			os.Exit(0)
		}
	}()
}

func recovery(roomID string) {
	err := recover()
	// no problem just shutdown
	if err == nil {
		return
	}

	// try to send that error to matrix and log, if available
	if bot != nil {
		bot.Error(id.RoomID(roomID), fatalmessage, err)
		return
	}
}
