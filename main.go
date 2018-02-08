package main

import (
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
)

func main() {
	// init logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("Loading config...")
	LoadConfig(sugar)
	wowBot := &Bot{SL: sugar}

	sugar.Info("Starting...")
	wowBot.Start()

	sugar.Info(http.ListenAndServe("localhost:6060", nil))
	<-make(chan struct{})
}
