package resources

import (
	
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/url"
	"strings"
)

type AudioStruct struct {
	Audio string
}

type IdStruct struct {
	Id string
}

func cooltown(w http.ResponseWriter, r *http.Request){
	var a AudioStruct
	err := json.NewDecoder(r.Body).Decode(&a); if err == nil{
		if a.Audio == "" {
			w.WriteHeader(400)
			return
		}
		posturl := "http://localhost:3001/search"
		body := []byte(`{"Audio": "`+ a.Audio +`"}`)

		r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body)); if err == nil{
			r.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			res, err := client.Do(r); if err == nil {
				if res.StatusCode == 400 {
					w.WriteHeader(400)
					return
				} else if res.StatusCode == 404 {
					w.WriteHeader(404)
					return
				} else if res.StatusCode == 500 {
					w.WriteHeader(500)
					return
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body); if err == nil{
					var mymap map[string]string
					json.Unmarshal(body, &mymap)
					var id = url.QueryEscape(strings.Replace(mymap["Id"], " ", "+", -1))
					posturl = "http://localhost:3000/tracks/" + id
					r, err := http.NewRequest("GET", posturl, nil); if err == nil {
						r.Header.Add("Content-Type", "application/json")
						client = &http.Client{}
						res, err := client.Do(r); if err == nil {
							if res.StatusCode == 404 {
								w.WriteHeader(404)
								return
							} else if res.StatusCode == 500 {
								w.WriteHeader(500)
								return
							}
							defer res.Body.Close()
							body, err := ioutil.ReadAll(res.Body); if err == nil {
								json.Unmarshal(body, &mymap)
								var audio = mymap["Audio"]
								var b = AudioStruct{Audio:audio}
								w.WriteHeader(200)
								json.NewEncoder(w).Encode(b)
								return
							}
						}
					}
					
				}

			}

		}
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(400)
	return
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/cooltown", cooltown).Methods("POST")

	return r
}
