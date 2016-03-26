package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// baby don't hurt me https://www.youtube.com/watch?v=xhrBDcQq2DM
const developerKey = "AIzaSyDJboG1vzku1lVFmtCsr6wbjKi7hyuw46E"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/msg", msgHandler)
	r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/s/index.html", 301)
}

func msgHandler(w http.ResponseWriter, r *http.Request) {

	q := "Akira - OST - Full Album"
	q = r.FormValue("body")
	originator := r.FormValue("originator")
	hash := fmt.Sprintf("%x", md5.Sum([]byte(originator)))
	userAvatarURL := "https://storage.googleapis.com/graphone-1337.appspot.com/avatars/" + hash + ".png"

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(q).
		Type("video").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		log.Println(item.Id.Kind)
		switch item.Id.Kind {
		case "youtube#video":
			// item.Id.VideoId
			post(item.Id.VideoId, item.Snippet.Title, userAvatarURL)
		}
	}
}

func post(id, title, userAvatar string) {
	log.Printf("- %s \n", id)
	entry := struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		UserAvatar string `json:"user_avatar"`
	}{ID: id, Title: title, UserAvatar: userAvatar}
	bs, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", "https://smsyoutube.firebaseio.com/rooms/one/messages.json", bytes.NewReader(bs))
	if err != nil {
		log.Println(err)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	contents, err := ioutil.ReadAll(r.Body)
	log.Println(string(contents))
}
