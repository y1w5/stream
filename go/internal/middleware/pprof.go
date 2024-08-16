package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/pprof"
)

const PprofCPUOutput = "cpu.prof"
const PprofMemOutput = "mem.prof"

func Pprof(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Create(PprofCPUOutput)
		if err != nil {
			logger.Error("could not create CPU profile", "err", err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Error("could not start CPU profile", "err", err)
		}
		defer pprof.StopCPUProfile()

		next.ServeHTTP(w, r)
	})
}
