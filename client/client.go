package client

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

type GameClient struct {
	Token string
}

func (gameClient *GameClient) InitGame() string {
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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	xAuthToken := res.Header.Get("X-Auth-Token")

	gameClient.Token = xAuthToken

	// fmt.Println(xAuthToken)

	return xAuthToken
}

type BoardStruct struct {
	Board []string `json:"board"`
}

func (gameClient *GameClient) Board() ([]string, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/board"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", gameClient.Token)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	defer res.Body.Close()

	var board BoardStruct
	err = json.NewDecoder(res.Body).Decode(&board)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return board.Board, err
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

// *StatusResponse
func (gameClient *GameClient) Status() (*GameStatus, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", gameClient.Token)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	defer res.Body.Close()

	var gameStatus GameStatus
	err = json.NewDecoder(res.Body).Decode(&gameStatus)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return &gameStatus, err
}

type FireStruct struct {
	Coord string `json:"coord"`
}

type FireResponse struct {
	Result string `json:"result"`
}

func (gameClient *GameClient) Fire(coord string) (*FireResponse, error) {
	posturl := "https://go-pjatk-server.fly.dev/api/game/fire"

	var fire FireStruct
	fire.Coord = coord
	jsonFire, _ := json.Marshal(fire)

	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(jsonFire))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Auth-Token", gameClient.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var fireResponse FireResponse
	err = json.NewDecoder(res.Body).Decode(&fireResponse)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return &fireResponse, err
}
