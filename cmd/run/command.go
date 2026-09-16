package run

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/matthiasharzer/go-pokemon-viewer/logging"
	"github.com/matthiasharzer/go-pokemon-viewer/queries/pokedex"
	"github.com/matthiasharzer/go-pokemon-viewer/ui"
	"github.com/matthiasharzer/go-pokemon-viewer/utils/httputils"
	"github.com/matthiasharzer/go-pokemon-viewer/views/pokemonview"
)

var httpPort int
var httpHost string
var fetchInterval = 24 * time.Hour

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "The HTTP server port to listen on")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "The HTTP server host (default: all interfaces)")
	Command.Flags().DurationVarP(&fetchInterval, "fetch-interval", "", fetchInterval, "The interval at which to fetch the pokémon data. If set to 0, will disable periodic updates, but will download pokémon date initially nonetheless.")
}

var Command = &cobra.Command{
	Use:          "run",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if fetchInterval < 0 {
			return fmt.Errorf("fetch interval cannot be negative")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pokemonView, err := pokemonview.New(context.Background(), fetchInterval)
		if err != nil {
			return fmt.Errorf("failed to create pokémon view: %w", err)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})
		mux.HandleFunc("GET /api/v1/pokedex", httputils.UseMiddleware(
			[]httputils.Middleware{
				httputils.GZIPMiddleware(),
				httputils.CacheMiddleware(24*time.Hour, func(r *http.Request) string {
					return pokemonView.Hash()
				}),
			},
			pokedex.Handler(pokemonView),
		))

		// Do not handle /api/* by the UI
		mux.Handle("GET /api/", http.NotFoundHandler())
		mux.Handle("GET /",
			httputils.UseMiddleware(
				[]httputils.Middleware{httputils.GZIPMiddleware()},
				httputils.HandleStaticSite(ui.Content),
			),
		)

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting go-pokemon-viewer-server", "host", httpHost, "port", httpPort)
		err = http.ListenAndServe(
			addr,
			mux,
		)

		return err
	},
}
