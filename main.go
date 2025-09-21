package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// flags
var (
	listenHostAndPort = flag.String("listenHostAndPort", "localhost:8080", "listen host and port")
	tcpHostAndPort    = flag.String("tcpHostAndPort", "localhost:31415", "tcp host and port")
	slogLevel         slog.Level
)

func parseFlags() {
	flag.TextVar(&slogLevel, "slogLevel", slog.LevelInfo, "slog level")

	flag.Parse()
}

func setupSlog() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slogLevel,
				},
			),
		),
	)

	slog.Info("setupSlog",
		"sloglevel", slogLevel,
	)
}

func buildInfoMap() map[string]string {
	buildInfoMap := make(map[string]string)

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		buildInfoMap["GoVersion"] = buildInfo.GoVersion
		for _, setting := range buildInfo.Settings {
			if strings.HasPrefix(setting.Key, "GO") ||
				strings.HasPrefix(setting.Key, "vcs") {
				buildInfoMap[setting.Key] = setting.Value
			}
		}
	}

	return buildInfoMap
}

func websocketServerHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		txID := uuid.New().String()

		txLogger := slog.Default().With(
			"txID", txID,
		)

		txLogger.Info("begin websocket handler",
			"method", r.Method,
			"headers", r.Header,
			"protocol", r.Proto,
			"url", r.URL.String(),
		)

		websocketConn, err := websocket.Accept(w, r, nil)
		if err != nil {
			txLogger.Warn("websocket.Accept error",
				"error", err,
			)
			return
		}

		defer websocketConn.CloseNow()

		tcpConn, err := net.DialTimeout("tcp", *tcpHostAndPort, 2*time.Second)
		if err != nil {
			txLogger.Warn("net.DialTimeout error",
				"error", err,
			)
			return
		}

		defer tcpConn.Close()

		wsNetConn := websocket.NetConn(context.Background(), websocketConn, websocket.MessageBinary)

		var proxyWaitGroup sync.WaitGroup

		proxyWaitGroup.Go(func() {
			defer wsNetConn.Close()
			defer tcpConn.Close()

			written, err := io.Copy(wsNetConn, tcpConn)

			txLogger.Info("after io.Copy(wsNetConn, tcpConn)",
				"written", written,
				"error", err,
			)
		})

		proxyWaitGroup.Go(func() {
			defer wsNetConn.Close()
			defer tcpConn.Close()

			written, err := io.Copy(tcpConn, wsNetConn)

			txLogger.Info("after io.Copy(tcpConn, wsNetConn)",
				"written", written,
				"error", err,
			)
		})

		proxyWaitGroup.Wait()

		txLogger.Info("end websocket handler")

	})
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic in main",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	parseFlags()

	setupSlog()

	slog.Info("begin main",
		"buildInfoMap", buildInfoMap(),
		"listenHostAndPort", *listenHostAndPort,
		"tcpHostAndPort", *tcpHostAndPort,
	)

	httpServer := &http.Server{
		Addr:         *listenHostAndPort,
		Handler:      websocketServerHandlerFunc(),
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	slog.Info("starting http server")

	err := httpServer.ListenAndServe()
	panic(fmt.Errorf("httpServer.ListenAndServe error: %w", err))
}
