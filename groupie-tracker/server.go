package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var artists []Artist
var homePage = template.Must(template.ParseFiles("templates/index.html"))
var artistPage = template.Must(template.ParseFiles("templates/artist.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        notFoundHandler(w, r)
        return
    }

    data, _, err := FetchArtists("") // fetch all artists
    if err != nil {
         internalErrorHandler(w, err)
        return
    }

    homePage.Execute(w, data)
}



func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowedHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	// fmt.Println(artists)
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	artist, _ := fetchData(id)
	if artist.Artist.ID == 0 {
		notFoundHandler(w, r)
		return
	}
	if err := artistPage.Execute(w, artist); err != nil {
		internalErrorHandler(w, err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound) // ставим код 404
    tmpl, _ := template.ParseFiles("404.html") // или "index.html" с сообщением
    tmpl.Execute(w, nil)
}


func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

// func badRequestHandler(w http.ResponseWriter, msg string) {
// 	http.Error(w, "400 bad request: "+msg, http.StatusBadRequest)
// }

// --------------------------------------------------------------------------------------/
func internalErrorHandler(w http.ResponseWriter, err error) {
    // Send plain 500 response to the client
    http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

