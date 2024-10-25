package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/handlers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/db"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/sb"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// embed public
// var FS embed.FS

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
	sbHost := os.Getenv("SUPABASE_URL")
	if len(sbHost) == 0 {
		log.Fatal("Neet supabase URL")
	}
	sbSecret := os.Getenv("SUPABASE_SECRET")
	if len(sbSecret) == 0 {
		log.Fatal("Need supabase token")
	}
	github_redirect_url := os.Getenv("GITHUB_AUTH_REDIRECT")
	if len(github_redirect_url) == 0 {
		log.Fatal("Need github callback redirect url")
	}

	bunDB, err := db.Init()
	if err != nil {
		log.Fatal("Cannot initialize the database", err)
	}

	psDB := db.SupabasePostgresql{
		Bun: bunDB,
	}

	sbClient := sb.InitDB(sbHost, sbSecret)
	// to change this when we pass the variable
	h := handlers.NewHandlers(*sbClient, github_redirect_url, &psDB)

	router := chi.NewMux()
	router.Use(h.IsLoggedIn)

	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	// router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/login", helpers.MakeHandler(h.HandleLoginIndex))
	router.Get("/login/provider/google", helpers.MakeHandler(h.HandleLoginGithub))
	router.Post("/login", helpers.MakeHandler(h.HandleLoginCreate))
	router.Get("/signup", helpers.MakeHandler(h.HandleSignupIndex))
	router.Post("/logout", helpers.MakeHandler(h.HandleLogoutCreate))
	router.Post("/signup", helpers.MakeHandler(h.HandleSignupCretate))
	router.Get("/auth/callback", helpers.MakeHandler(h.HandleAuthCallback))
	router.Get("/auth/v1/callback", helpers.MakeHandler(h.HandleAuthCallback))
	router.Get("/account/setup", helpers.MakeHandler(h.HandleAccountSetupIndex))
	router.Post("/account/setup", helpers.MakeHandler(h.HandleAccountSetupCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(h.WithAccountSetup)
		auth.Get("/", helpers.MakeHandler(h.HandleHomeIndex))
		auth.Get("/settings", helpers.MakeHandler(h.HandleSettingsIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application is running", "port", port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router))
}

func getEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
