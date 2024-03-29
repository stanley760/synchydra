package main

import (
	"fmt"
	"go.uber.org/zap"
	"synchydra/pkg/config"
	"synchydra/pkg/http"
	"synchydra/pkg/log"
)

func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)

	servers, cleanup, err := newApp(conf, logger)
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", "http://localhost:"+conf.GetString("http.port")))

	//servers.
	http.Run(servers.ServerHTTP, fmt.Sprintf(":%d", conf.GetInt("http.port")))
	defer cleanup()

}
