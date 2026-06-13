package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
)

type Options struct {
	Addr string
}

func Run(opts Options) error {
	if opts.Addr == "" {
		opts.Addr = "127.0.0.1:18181"
	}

	cfgDir, err := localconfig.ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		return err
	}
	pidPath := filepath.Join(cfgDir, "daemon.pid")
	_ = os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", os.Getpid())), 0600)
	defer os.Remove(pidPath)

	startedAt := time.Now()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":             true,
			"build_id":       buildinfo.BuildID,
			"build_channel":  buildinfo.BuildChannel,
			"client_version": buildinfo.Version,
			"uptime_seconds": int(time.Since(startedAt).Seconds()),
		})
	})

	ln, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	_ = logs.Append("daemon-start", map[string]string{"addr": opts.Addr})

	server := &http.Server{Handler: mux}
	err = server.Serve(ln)
	_ = logs.Append("daemon-stop", map[string]string{
		"addr": opts.Addr,
	})
	return err
}
