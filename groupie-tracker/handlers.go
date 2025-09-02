package main

import (
	"fmt"
	"net/http"
)

func Execute(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        notFoundHandler(w, r)
        fmt.Println(r.URL.Path)
        return
    }

    targetName := r.URL.Query().Get("artist") // get search query
    data, _,err := FetchArtists(targetName)
    if err != nil {
        http.Error(w, "We cannot get artists", http.StatusInternalServerError)
        return
    }

    homePage.Execute(w, data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    artistName := r.URL.Query().Get("artist")
    if artistName == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    _, found, err := FetchArtists(artistName)
    if err != nil {
        http.Error(w, "500 internal server error", http.StatusInternalServerError)
        return
    }

    if found != nil {
        http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", found.ID), http.StatusSeeOther)
        return
    }

    http.Error(w, "404 page not found", http.StatusNotFound)
}




func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "No artist ID provided", http.StatusBadRequest)
		return
	}

	artist, err := http.Get(id)
	if err != nil {
		http.Error(w, "Cannot fetch artist", http.StatusInternalServerError)
		return
	}

	artistPage.Execute(w, artist)
}
