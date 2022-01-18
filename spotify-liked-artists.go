package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

const AUTH_URL string = "https://accounts.spotify.com/api/token"
const LIKES_URL string = "https://api.spotify.com/v1/me/tracks"

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Response struct {
	Href     string `json:"href"`
	Items    []Song `json:"items"`
	Limit    int32  `json:"limit"`
	Next     string `json:"next"`
	Offset   int32  `json:"offset"`
	Previous string `json:"previous"`
	Total    int32  `json:"total"`
}

type Song struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

type Track struct {
	Album   map[string]interface{} `json:"Album"`
	Artists []Artist               `json:"Artists"`
}

type Artist struct {
	ExternalUrls map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Uri          string            `json:"uri"`
}

func main() {
	loadEnvFile("secrets.env")

	// client_id := os.Getenv("CLIENT_ID")
	// client_secret := os.Getenv("CLIENT_SECRET")

	// token := getToken(client_id, client_secret)
	token := os.Getenv("TOKEN")

	// fmt.Println(token)

	tracks := getTracks(token)

	artists := getArtists(tracks)

	pretty, err := json.MarshalIndent(artists, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(pretty))
}

func getToken(client_id, client_secret string) string {

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", client_id)
	params.Add("client_secret", client_secret)

	resp, err := http.PostForm(AUTH_URL, params)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response_map Token

	err = json.Unmarshal(body, &response_map)
	if err != nil {
		log.Fatal(err)
	}

	return response_map.AccessToken
}

func getTracks(token string) []Song {
	client := &http.Client{}

	req, err := http.NewRequest("GET", LIKES_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// body_string := string(body)
	// fmt.Println(body_string)

	var response_map Response

	err = json.Unmarshal(body, &response_map)
	if err != nil {
		log.Fatal(err)
	}

	return response_map.Items
}

func getArtists(songs []Song) []string {

	artist_list := []string{}
	artist_map := map[string]bool{}

	for _, song := range songs {
		artists := song.Track.Artists

		for _, artist := range artists {
			artist_map[artist.Name] = true
		}
	}

	for name := range artist_map {
		artist_list = append(artist_list, name)
	}

	sort.Strings(artist_list)

	return artist_list
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	} else {
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		cur_string := scanner.Text()
		env_var := strings.Split(cur_string, "=")
		os.Setenv(env_var[0], env_var[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
