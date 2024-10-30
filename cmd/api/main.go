// Luigi Acuna
// CMPS4191 Test 3 Advanced Web Dev
// October 30 2024
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port       int    //port number to access signin page
	enviroment string //enviroment the signin page will be on
}

type applicationDependencies struct {
	config serverConfig
	logger *slog.Logger //look more into this later
}

func main() {
	var settings serverConfig

	//Settings ports and enviroment info
	flag.IntVar(&settings.port, "port", 4000, "Server Port")
	flag.StringVar(&settings.enviroment, "env", "development", "Enviroment(development|staging|)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	appInstance := &applicationDependencies{
		config: settings,
		logger: logger,
	}

	//api server info
	apiServer := http.Server{
		Addr:         fmt.Sprintf(":%d", settings.port),
		Handler:      appInstance.routes(), //this one too
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "address", apiServer.Addr, "enviroment", settings.enviroment)

	err := apiServer.ListenAndServe()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
