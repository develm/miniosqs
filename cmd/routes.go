package cmd

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) Routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.notifySqs)
	mux.HandleFunc("/bucket/setup", app.setupMinioBucket)

	return standardMiddleware.Then(mux)
}
