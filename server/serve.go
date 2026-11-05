package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	apischema "github.com/2manyvcos/paranal/api"
	"github.com/2manyvcos/paranal/client"
	"github.com/2manyvcos/paranal/server/api"
	"github.com/2manyvcos/paranal/server/config"
)

const API_PATH = "/api"

func Serve(ctx context.Context, cfg *config.Config) error {
	app := api.NewApp(cfg)

	mux := http.NewServeMux()

	apiHandler := api.New()
	mux.Handle(API_PATH+"/", http.StripPrefix(API_PATH, api.OmitTrailingSlash(api.WithApp(app, api.WithCORS(apiHandler)))))

	mux.Handle("/", http.FileServer(api.FileRewrite(api.FileTemplates(client.ClientFiles, map[string]any{
		"appName": cfg.AppName,
		"api":     cfg.PublicAPI,
	}, "index.html", "manifest.json"), "index.html")))

	mux.Handle("/api/reference/", http.FileServer(api.FileTemplates(client.ClientFiles, map[string]any{
		"api": cfg.PublicAPI,
	}, "api/reference/index.html")))

	apiSchemaFiles, _ := fs.Glob(apischema.SchemaFiles, "*.yaml")
	apiSchemaServer := api.WithApp(app, api.WithCORS(http.FileServer(api.FileTemplates(apischema.SchemaFiles, map[string]any{
		"api": cfg.PublicAPI,
	}, apiSchemaFiles...))))
	mux.Handle("/api/schema/", http.StripPrefix("/api/schema/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if ext := filepath.Ext(req.URL.Path); ext == ".yaml" || ext == ".yml" {
			res.Header().Set("Content-Type", "text/yaml")
		}
		apiSchemaServer.ServeHTTP(res, req)
	})))

	srv := &http.Server{
		Addr:        fmt.Sprintf("%v:%v", cfg.Server.Address, cfg.Server.Port),
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	var serve func() error
	switch strings.ToLower(cfg.Server.Protocol) {
	case "http":
		serve = func() error {
			log.Printf("listening at http://%v\n", srv.Addr)
			return srv.ListenAndServe()
		}

	case "https":
		serve = func() error {
			log.Printf("listening at https://%v\n", srv.Addr)
			return srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile)
		}

	default:
		return fmt.Errorf("unsupported protocol %q", cfg.Server.Protocol)
	}

	serveErr := make(chan error, 1)
	go func() {
		err := serve()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}

		close(serveErr)
	}()

	select {
	case err := <-serveErr:
		return err

	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gracefully shutting down server: %w", err)
	}

	return nil
}
