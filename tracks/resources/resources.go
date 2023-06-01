package resources

import (
	"encoding/json"
	"net/http"
	"tracks/repository"

	"github.com/gorilla/mux"
)

func updateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var t repository.Track
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if id == t.Id {
			if n := repository.Update(t); n > 0 {
				w.WriteHeader(204) /* No Content */
			} else if n := repository.Insert(t); n > 0 {
				w.WriteHeader(201) /* Created */
			} else {
				w.WriteHeader(500) /* Internal Server Error */
			}
		} else {
			w.WriteHeader(400) /* Bad Request */
		}
	} else {
		w.WriteHeader(400) /* Bad Request */
	}

}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if t, n := repository.Read(id); n > 0 {
		d := repository.Track{Id: t.Id, Audio: t.Audio}
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func listTracks(w http.ResponseWriter, r *http.Request) {
	if tracklist, n := repository.List(); n > 0 {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(tracklist)
	} else {
		w.WriteHeader(500)
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if n := repository.Delete(id); n > 0 {
		w.WriteHeader(204)
	} else if n == 0 {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(500)
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/tracks/{id}", updateTrack).Methods("PUT")
	r.HandleFunc("/tracks/{id}", readTrack).Methods("GET")
	r.HandleFunc("/tracks", listTracks).Methods("GET")
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")

	return r
}
