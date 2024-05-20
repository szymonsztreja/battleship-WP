package temp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// type Post struct {
// 	Coords
// 	Desc
// 	Nick
// 	Target_nick
// }

type HttpGameClient struct {
	Client     *http.Client
	XAuthToken string
}

type BoardStruct struct {
	Board []string `json:"board"`
}

type GameStatus struct {
	Nick           string   `json:"nick"`
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Opponent       string   `json:"opponent"`
	OppShots       []string `json:"opp_shots"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

type FireStruct struct {
	Coord string `json:"coord"`
}

type FireResponse struct {
	Result string `json:"result"`
}

type PlayersDescription struct {
	Desc     string `json:"desc"`
	Nick     string `json:"nick"`
	OppDesc  string `json:"opp_desc"`
	Opponent string `json:"opponent"`
}

func (httpClient *HttpGameClient) InitGame() (int, error) {
	posturl := "https://go-pjatk-server.fly.dev/api/game"

	body := []byte(`{
        "coords": [
    "A1",
    "A3",
    "B9",
    "C7",
    "D1",
    "D2",
    "D3",
    "D4",
    "D7",
    "E7",
    "F1",
    "F2",
    "F3",
    "F5",
    "G5",
    "G8",
    "G9",
    "I4",
    "J4",
    "J8"
  ],
  "desc": "My first game",
  "nick": "AAAAAAAAA",
  "target_nick": "",
  "wpbot": true
    }`)

	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	xAuthXAuthToken := res.Header.Get("X-Auth-Token")

	httpClient.XAuthToken = xAuthXAuthToken

	return res.StatusCode, err
}

func (httpClient *HttpGameClient) Board() ([]string, int, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/board"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.Client.Do(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	defer res.Body.Close()

	var b BoardStruct
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return b.Board, res.StatusCode, err
}

func (httpClient *HttpGameClient) Status() (*GameStatus, int, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.Client.Do(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}
	defer res.Body.Close()

	var gameStatus GameStatus
	err = json.NewDecoder(res.Body).Decode(&gameStatus)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}
	return &gameStatus, res.StatusCode, err
}

func (httpClient *HttpGameClient) Fire(coord string) (*FireResponse, int, error) {
	posturl := "https://go-pjatk-server.fly.dev/api/game/fire"

	var fire FireStruct
	fire.Coord = coord
	jsonFire, _ := json.Marshal(fire)

	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(jsonFire))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Client.Do(req)
	if err != nil {
		fmt.Printf("error getting http description request: %s\n", err)
	}

	defer res.Body.Close()

	var fireResponse FireResponse
	err = json.NewDecoder(res.Body).Decode(&fireResponse)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return &fireResponse, res.StatusCode, err
}

func (httpClient *HttpGameClient) GetPlayersDescription() (*PlayersDescription, int, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/desc"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.Client.Do(req)

	if err != nil {
		fmt.Printf("error getting http description request: %s\n", err)
	}
	defer res.Body.Close()

	var desc PlayersDescription
	err = json.NewDecoder(res.Body).Decode(&desc)
	if err != nil {
		fmt.Printf("error decoding http description request: %s\n", err)
	}
	return &desc, res.StatusCode, err
}
