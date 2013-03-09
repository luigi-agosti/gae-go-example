package imageapi

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
    "appengine"
    "appengine/datastore"
)

type Image struct {
	Key string `json:"key"`
    Url string `json:"url"`
    ThumbUrl string `json:"thumbnailUrl"`
    Description string `json:"description"`
    Link string `json:"link"`
    Service string `json:"service"`
}

func init() {
    http.HandleFunc("/image", handleImage)
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		imageGet(w, r)
	} else if r.Method == "POST" {
		imagePost(w, r)
	}
}

func imageGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    q := datastore.NewQuery("Image").Limit(10)
    images := make([]Image, 0, 10)
    if _, err := q.GetAll(c, &images); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    bytes, _ := json.Marshal(images)
	w.Header().Add("content-type", "application/json")
	w.Write(bytes)
}

func imagePost(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var m Image
	json.Unmarshal(body, &m)
	w.WriteHeader(http.StatusOK)
	c := appengine.NewContext(r)
 	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Image", nil), &m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	bytes, _ := json.Marshal(m)
	w.Header().Add("content-type", "application/json")
	w.Write(bytes)
}
