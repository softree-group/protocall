package interfaces

import "github.com/gorilla/mux"

func NewRouter(app *Application) *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/protocols", app.MakeProtocol).Methods("POST")
	mux.HandleFunc("/records", app.ProcessRecord).Methods("POST")
	mux.HandleFunc("/records", app.GetJobStatus).Methods("GET").Queries("{id:[0-9]+}")
	return mux
}
