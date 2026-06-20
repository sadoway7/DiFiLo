package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"difilo/internal/app"
	"difilo/internal/auth"
	"difilo/internal/content"
	"difilo/internal/db"
)

func main() {
	var (
		mirror   = flag.String("mirror", "./mirror", "path to the mirror directory")
		host     = flag.String("host", "127.0.0.1", "interface to bind (use 0.0.0.0 to expose)")
		port     = flag.Int("port", 8000, "port to serve on")
		reimport = flag.Bool("reimport", false, "force re-import of content from mirror/md/")
		dbPath   = flag.String("db", "", "path to the SQLite database (default: <mirror>/difilo.db)")
	)
	flag.Parse()

	app.InitLogger(false)

	dbFile := *dbPath
	if dbFile == "" {
		dbFile = filepath.Join(*mirror, "difilo.db")
	}

	database, err := db.OpenDB(dbFile)
	if err != nil {
		die("opening database: %v", err)
	}

	if *reimport || !database.ContentImported() {
		if *reimport {
			slog.Info("re-importing content (forced)")
		} else {
			slog.Info("importing content (first run)")
		}
		stats, err := content.ImportContent(database, *mirror)
		if err != nil {
			die("importing content: %v", err)
		}
		slog.Info("import complete",
			"pages", stats.Pages,
			"images", stats.Images,
			"links", stats.Links,
			"duration", stats.Duration,
		)
	}

	auth.InitSecret(nil)

	srv := app.New(app.DefaultConfig(), database, *mirror)
	srv.BuildHeroImages()
	srv.BuildAliases()

	mux := http.NewServeMux()
	mux.Handle("/", srv.Handler())
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", addr, "mirror", *mirror)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			die("server: %v", err)
		}
	}()

	<-stop
	slog.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
	slog.Info("server stopped")
}

func die(format string, a ...any) {
	slog.Error("fatal", "error", fmt.Sprintf(format, a...))
	os.Exit(1)
}
