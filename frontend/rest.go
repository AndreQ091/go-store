package frontend

import (
	"errors"
	"go-store/core"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type restFrontend struct {
	store *core.KeyValueStore
}

func (f *restFrontend) Start(store *core.KeyValueStore) error {
	f.store = store
	r := mux.NewRouter()
	r.HandleFunc("/key/{key}", f.keyValueGetHandler).Methods("GET")
	r.HandleFunc("/key/{key}", f.keyValuePutHandler).Methods("PUT")

	return http.ListenAndServe(":8000", r)
}

func (f *restFrontend) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := f.store.Get(key)

	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func (f *restFrontend) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = f.store.Put(key, string(value))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
