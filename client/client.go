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
	token string
}

func InitGame(gameClient *GameClient) string {
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

	gameClient.token = xAuthToken

	fmt.Println(xAuthToken)

	return xAuthToken
}

type BoardStruct struct {
	Board []string `json:"board"`
}

func Board(gameClient *GameClient) ([]string, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/board"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", gameClient.token)
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
	GameStatus     string `json:"board"`
	LastGameStatus string `json:"last_game_status"`
	Nick           string `json:"nick"`
	OppShots       string `json:"opp_shots"`
	Opponent       string `json:"opponent"`
	ShouldFire     bool   `json:"should_fire"`
	Timer          int    `json:"timer"`
}

// *StatusResponse
func Status(gameClient *GameClient) (int, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/board"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", gameClient.token)
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

	return res.StatusCode, err
}

type FireStruct struct {
	Coord string `json:"coord"`
}

func Fire(gameClient *GameClient, coord string) (string, error) {
	posturl := "https://go-pjatk-server.fly.dev/api/game/fire"

	var fire FireStruct
	fire.Coord = coord
	jsonFire, _ := json.Marshal(fire)

	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(jsonFire))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Auth-Token", gameClient.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	return
}
