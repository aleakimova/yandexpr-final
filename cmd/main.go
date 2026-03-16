package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/aleakimova/yandexpr-final/internal/db"
	"github.com/aleakimova/yandexpr-final/pkg/api"
)

const defaultPort = 7540
const defaultDBFile = "scheduler.db"

func main() {
	logLevel := slog.LevelInfo
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		if err := logLevel.UnmarshalText([]byte(envLevel)); err != nil {
			slog.Warn("invalid LOG_LEVEL, using INFO", "value", envLevel)
		}
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// Read cli args
	args := os.Args[1:]
	if len(args) < 1 {
		slog.Error("no html dir specified")
		os.Exit(1)
	}

	// Override default DB file with TODO_DB
	dbfile := defaultDBFile
	envFile := os.Getenv("TODO_DB")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	slog.Info("opening database", "file", dbfile)

	// Start SQLite DB
	_, err := db.Start(dbfile)
	if err != nil {
		slog.Error("failed to initialize database", "file", dbfile, "error", err)
		panic(err)
	}

	// Override default port with TODO_PORT
	port := defaultPort
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		} else {
			slog.Warn("invalid TODO_PORT value, using default", "value", envPort, "default", defaultPort)
		}
	}

	http.Handle("/", http.FileServer(http.Dir(args[0])))
	api.Init()

	addr := fmt.Sprintf(":%d", port)
	slog.Info("server starting", "addr", addr, "static_dir", args[0])
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server stopped", "error", err)
		panic(err)
	}
}
