package services

import (
	"os"
	"os/signal"
)

// InterruptCh returns a channel that listens for OS interrupt signals.
func InterruptCh() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	return stop
}