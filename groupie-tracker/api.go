package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// region types
// Все Артисты 
type PageData struct {
    Artists []Artist
    Count   int
}

// Artist хранит данные об артисте
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	FirstAlbum   string   `json:"firstAlbum"`
	CreationDate int      `json:"creationDate"`
}


type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// ----------------------------------------------------------------------/
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// ---------------------------------------------------------------------------/
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistDetails struct {
	*Artist
	*Dates
	*Relation
	*Locations
}
// endregion types

// Загружаем артистов с API
func FetchArtists(targetName string) (PageData, *Artist, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
    if err != nil {
        return PageData{}, nil, err
    }
    defer resp.Body.Close()

    var artists []Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return PageData{}, nil, err
    }

    var found *Artist
    if targetName != "" {
		targetLower := strings.ToLower(targetName)
        for i, a := range artists {
              if strings.ToLower(a.Name) == targetLower {
                found = &artists[i]
                break
            }
        }
    }

    data := PageData{
        Artists: artists,
        Count:   len(artists),
    }

    return data, found, nil
}




func FetchArtist(id int) (*Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists *Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

// ----------------------------------------------------------------------------/
func FetchJSon(url string, artist interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(artist)
}

func fetchData(id int) (*ArtistDetails, error) {
	/// id  подставляем автоматически
	artist, artistErr := FetchArtist(id)
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(id)
	datesURL := "https://groupietrackers.herokuapp.com/api/dates/" + strconv.Itoa(id)
	relationURL := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)
	locations := &Locations{}
	dates := &Dates{}
	relation := &Relation{}
	locErr := FetchJSon(locationsURL, locations)
	datesErr := FetchJSon(datesURL, dates)
	relationsErr := FetchJSon(relationURL, relation)

	return &ArtistDetails{
		artist, dates, relation, locations,
	}, errors.Join(artistErr, locErr, datesErr, relationsErr)
}


