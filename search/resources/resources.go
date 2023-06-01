package resources

import (
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
)

type AudioStruct struct {
	Audio string
}

type IdStruct struct {
	Id string
}

func searchTrack(w http.ResponseWriter, r *http.Request) {
	var a AudioStruct
	err := json.NewDecoder(r.Body).Decode(&a); if err == nil{
		if a.Audio == "" {
			w.WriteHeader(400)
			return
		}
		posturl := "https://api.audd.io"
		body := []byte(`{"audio": "`+ a.Audio +`", "api_token": "test"}`)

		r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body)); if err == nil{
			r.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			res, err := client.Do(r); if err == nil {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body); if err == nil{
					var mymap map[string]string
					var response map[string]map[string]string
					json.Unmarshal(body, &response)
					json.Unmarshal(body, &mymap)
					if mymap["status"] == "success"{
						if response["result"] != nil {
							var t = IdStruct{Id: response["result"]["title"]}
							w.WriteHeader(200)
							json.NewEncoder(w).Encode(t)
							return	
						} else {
							w.WriteHeader(404)
							return
						}
					}
				}		
			}
		}
		w.WriteHeader(500)
		return

	}
	w.WriteHeader(400)	 
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/search", searchTrack).Methods("POST")

	return r
}
