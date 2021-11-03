package interfaces

import "github.com/gorilla/mux"

func NewRouter(app *Application) *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/protocols", app.create).Methods("POST")
	mux.HandleFunc("/records", app.translate).Methods("POST")
	return mux
}
