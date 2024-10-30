package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func defaultSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // kill -15
		syscall.SIGQUIT, // Ctrl+\
	}
}

func Context(ctx context.Context, sig ...os.Signal) context.Context {
	if len(sig) == 0 {
		sig = defaultSignals()
	}

	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, sig...)

	go func() {
		defer func() {
			signal.Stop(sigChan)
			cancel()
		}()

		select {
		case sig := <-sigChan:
			zerolog.Ctx(ctx).Info().Str("signal", sig.String()).Send()

			return

		case <-ctx.Done():
			zerolog.Ctx(ctx).Info().Err(ctx.Err()).Msg("context cancelled")

			return
		}
	}()

	return ctx
}
