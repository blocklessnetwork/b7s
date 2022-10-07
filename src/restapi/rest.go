package restapi

import (
	"context"
	"net/http"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func handleWeb(w http.ResponseWriter, r *http.Request) {
	// params
	w.Write([]byte("ok"))
}

func startServer(ctx context.Context) {
	var config = ctx.Value("config").(models.Config)
	// router for api
	myRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api/v1").Subrouter()
	myRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// router declaration
	myRouter.HandleFunc("/", handleRootRequest)

	myRouter.HandleFunc("/function/request", handleRequestExecute).Methods("POST")
	myRouter.HandleFunc("/function/install", handleInstallFunction).Methods("POST")
	myRouter.HandleFunc("/function/result", handleGetExecuteResponse).Methods("POST")

	log.Info(http.ListenAndServe(":"+config.Rest.Port, myRouter))
}

func Start(ctx context.Context) {
	var config = ctx.Value("config").(models.Config)

	log.WithFields(log.Fields{
		"port":    config.Rest.Port,
		"address": config.Rest.IP,
	}).Info("starting rest server")

	go startServer(ctx)
}
