package utils

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func GraceFullyShutDown(rootCtx context.Context, cleanup func(ctx context.Context)) {
	shutdown, stop := signal.NotifyContext(rootCtx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// waiting to have a disrupt
	<-shutdown.Done()

	log.Println("gracefully shutting down in progress...")

	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	cleanup(ctx)

	log.Println("shutdown complete")
}
