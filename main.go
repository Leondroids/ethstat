package main

import (
	"github.com/Leondroids/ethstat/app"
	"log"
	"github.com/rs/cors"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/Leondroids/gox"
)

func main() {
	context, err := app.InitApp()

	if err != nil {
		panic(err)
	}

	log.Println(context.Config)

	err = app.StartSyncingBlockDate(context)

	if err != nil {
		panic(err)
	}

	startServer(context)
}

func startServer(context *app.Context) error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	log.Printf("Starting Server at %v%v", host, context.Config.Port)
	handler := cors.Default().Handler(router(context))
	log.Fatal(http.ListenAndServe(context.Config.Port, handler))

	return nil
}

func router(context *app.Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	commonHandlers := alice.New(gox.LoggingHandler)

	router.Handle("/healthcheck", commonHandlers.ThenFunc(gox.HealthcheckHandler)).Methods("GET")

	dateByBlock := app.NewDateByBlock(context)
	router.Handle("/datebyblock", commonHandlers.ThenFunc(dateByBlock.GetDateByBlock)).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	return router
}
