package cmd

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/Mycelium/internal/encre"
	"github.com/FacileStudio/Mycelium/internal/env"
	"github.com/FacileStudio/Mycelium/internal/server"
	"github.com/FacileStudio/tronc/logger"
	"github.com/FacileStudio/tronc/spa"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Run sync server with web UI API",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := env.Load()
		if err != nil {
			logger.New(logger.Config{}).Error("failed to load config", slog.Any("error", err))
			os.Exit(1)
		}
		if dataDir, _ := cmd.Flags().GetString("data"); dataDir != "" {
			cfg.DataDir = dataDir
		}
		if cmd.Flags().Changed("port") {
			cfg.Port = servePort
		}

		var journalClient *journal.Client
		appLogger := logger.New(logger.Config{
			Level: cfg.LogLevel,
			Wrap: func(handler slog.Handler) slog.Handler {
				if cfg.JournalURL == "" || cfg.JournalToken == "" {
					return handler
				}
				journalClient = journal.New(journal.Config{URL: cfg.JournalURL, Token: cfg.JournalToken})
				return journal.NewHandler(journalClient, handler)
			},
		})
		if journalClient != nil {
			defer journalClient.Close()
		}

		srv := server.New(cfg.DataDir, cfg.Password)
		srv.Version = version
		srv.Edge = server.EdgeConfig{
			TrustedProxies: cfg.TrustedProxies,
			CDNProxies:     cfg.CDNProxies,
			CDNHeader:      cfg.CDNHeader,
		}
		srv.OIDC = cfg.OIDC
		srv.CORSAllowedOrigins = cfg.CORSAllowedOrigins
		srv.JournalBrowserURL = cfg.JournalBrowserURL
		srv.JournalBrowserKey = cfg.JournalBrowserKey
		srv.Log = appLogger

		if err := srv.EnableMCP(); err != nil {
			appLogger.Error("invalid MCP configuration", slog.Any("error", err))
			os.Exit(1)
		}

		emitter := server.NewEmitter(srv)
		go emitter.Run(cmd.Context())
		go srv.RunRegistreSync(cmd.Context())
		go srv.RunDeletions(cmd.Context())

		clientDir := spa.DirFromEnv()
		if hashes, err := spa.InlineScriptHashes(clientDir); err == nil {
			srv.Edge.ScriptSources = hashes
		}
		router := srv.Handler()
		if spa.Available(clientDir) {
			router.Handle("/*", spa.Handler(spa.Config{Dir: clientDir}))
			appLogger.Info("serving client", slog.String("dir", clientDir))
		}

		if cfg.Password == "" && cfg.OIDC == nil {
			appLogger.Warn("no authentication configured, every request is served as admin (set PASSWORD or OIDC_ISSUER)")
		}
		if !encre.Available() {
			appLogger.Error("encre is not installed: artifacts written through the API are stored without validation (set ENCRE_BIN or put encre on PATH)")
		}

		addr := ":" + strconv.Itoa(cfg.Port)
		appLogger.Info("server starting",
			slog.String("addr", addr),
			slog.String("data_dir", cfg.DataDir),
			slog.String("env", string(cfg.AppEnv)))
		return http.ListenAndServe(addr, router)
	},
}

func newServeCmd() *cobra.Command {
	serveCmd.Flags().IntVar(&servePort, "port", env.DefaultPort, "Port to listen on (default: $PORT)")
	serveCmd.Flags().String("data", "", "Data directory (default: $DATA_DIR, else ${XDG_DATA_HOME:-~/.local/share}/mycelium)")
	rootCmd.AddCommand(serveCmd)
	return serveCmd
}
