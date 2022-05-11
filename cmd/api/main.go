package main

import (
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hojabri/geolocation/pkg/api"
	"github.com/hojabri/geolocation/pkg/config"
	"github.com/hojabri/geolocation/pkg/maxmind"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "log level [trace|debug|info|warn|error]")
	flag.Parse()

	// log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch *logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	// pretty log - should be disabled in production
	log = log.Output(zerolog.NewConsoleWriter())

	// configuration setup
	config.Setup(log)

	maxmindServer := maxmind.New(&log, config.Configuration.GetString("MAXMIND_LICENSE_KEY"))

	// Download Maxmind GeoLite DB file when app starts
	if err := maxmindServer.DownloadDB(); err != nil {
		log.Err(err).Msg("Could not download the Geo DB file")
	}

	// Schedule download Maxmind DB file weekly
	if err := maxmindServer.RunDownloadScheduler(); err != nil {
		log.Err(err).Msg("Could not download the Geo DB file")
	}

	// Connect to GeoLite2 City db
	err := maxmindServer.OpenDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Can not connect to GeoLite2 City db")
	}
	defer func() {
		err = maxmindServer.CloseDB()
		if err != nil {
			log.Fatal().Err(err).Msg("Can not close the GeoLite2 City db")
		}
	}()

	// new api server
	s := api.New(&log, maxmindServer.DB)

	// create go-chi client
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer, middleware.RealIP, middleware.RequestID, middleware.Compress(5))
	r.Get("/health", Health)
	r.Mount("/api/v1", s.Routes())

	// setup ReDoc
	fs := http.FileServer(http.Dir("openapi"))
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	swaggerUrl := "/openapi/openapi.yml"
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err = tmpl.Execute(w, map[string]string{
			"url": swaggerUrl,
		})
		if err != nil {
			err = api.SendHTTPError(w, "failed to generate API documentation", http.StatusBadRequest)
			if err != nil {
				log.Error().Err(err).Msg("can't send error message")
			}
			return
		}
	})

	srv := &http.Server{
		Addr:              ":3000",
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		Handler:           r,
	}

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info().Msg("received signal, stopping server..")
		_ = srv.Shutdown(context.Background())
	}()

	log.Info().Str("addr", srv.Addr).Msg("starting api server")
	if err := srv.ListenAndServe(); err != nil {
		log.Info().Err(err).Msg("server stopped")
	}

}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
