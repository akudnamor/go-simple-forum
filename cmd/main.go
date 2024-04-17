package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go-simple-forum/internal/config"
	"go-simple-forum/internal/http-server/handler"
	"go-simple-forum/internal/lib/logger/handlers/slogpretty"
	"go-simple-forum/internal/lib/logger/sl"
	"go-simple-forum/internal/storage"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	/*
		TODO:
			передача в storage log`а
	*/

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Start!", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	router := chi.NewRouter()
	st, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	t, err := template.ParseFiles("internal/template/html/index.html", "internal/template/html/header.html", "internal/template/html/signup.html", "internal/template/html/auth.html")
	if err != nil {
		log.Error("failed to init template", sl.Err(err))
	}

	router.Handle("/internal/*", http.StripPrefix("/internal", http.FileServer(http.Dir("internal"))))

	router.Get("/", handler.IndexPage(log, st, t))

	router.Get("/signup", handler.SignUpPage(log, st, t))
	router.Get("/auth", handler.AuthPage(log, st, t))

	router.Post("/signup", handler.SignUp(log, st))
	router.Post("/auth", handler.Auth(log, st))
	router.HandleFunc("/logout", handler.Logout(log, st))

	/*
		TODO:
			/profile/{profile_id}
			/forum/{forum_id}
	*/

	s := http.Server{
		Addr:         cfg.Address,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		Handler:      router,
	}
	log.Info("Server is listening on", slog.String("address", cfg.Address))
	if err = s.ListenAndServe(); err != nil {
		log.Error("failed to start", sl.Err(err))
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = setupPrettySlog()
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	prettyHandler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(prettyHandler)
}
