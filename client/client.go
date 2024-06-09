package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GameData struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	Wpbot      bool     `json:"wpbot"`
}

const retry = 5

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

type MessageResponse struct {
	Message string `json:"message"`
}

type PlayerStatus struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`
}

type Top10Players struct {
	Stats []PlayerStats `json:"stats"`
}

type PlayerStats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

func (httpClient *HttpGameClient) makeRequest(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	for i := 0; i < retry; i++ {
		res, err = httpClient.Client.Do(req)
		if err != nil {
			fmt.Printf("Error making http request: %s\n", err)
			continue
		}
		if res.StatusCode != 200 {
			handleResponseCode(res.StatusCode, req)
			time.Sleep(350 * time.Millisecond)
		} else {
			break
		}

	}
	return res, err
}

func (httpClient *HttpGameClient) InitGame(gd GameData) {
	posturl := "https://go-pjatk-server.fly.dev/api/game"

	body, err := json.Marshal(gd)
	if err != nil {
		fmt.Printf("error marshaling game data: %v", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.makeRequest(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	xAuthXAuthToken := res.Header.Get("X-Auth-Token")

	httpClient.XAuthToken = xAuthXAuthToken
}

func (httpClient *HttpGameClient) Board() ([]string, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/board"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	defer res.Body.Close()

	var b BoardStruct
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return b.Board, err
}

// *StatusResponse
func (httpClient *HttpGameClient) Status() (*GameStatus, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.makeRequest(req)

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

func (httpClient *HttpGameClient) Fire(coord string) (*FireResponse, error) {
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

	res, err := httpClient.makeRequest(req)
	if err != nil {
		fmt.Printf("error getting http description request: %s\n", err)
	}

	defer res.Body.Close()

	var fireResponse FireResponse
	err = json.NewDecoder(res.Body).Decode(&fireResponse)
	if err != nil {
		fmt.Printf("error decoding http request: %s\n", err)
	}

	return &fireResponse, err
}

func (httpClient *HttpGameClient) GetPlayersDescription() (*PlayersDescription, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/desc"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error getting http description request: %s\n", err)
	}
	defer res.Body.Close()

	var desc PlayersDescription
	err = json.NewDecoder(res.Body).Decode(&desc)
	if err != nil {
		fmt.Printf("error decoding http description request: %s\n", err)
	}
	return &desc, err
}

func (httpClient *HttpGameClient) RefreshSession() (*MessageResponse, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/refresh"

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error getting refresh request: %s\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		msgRes := &MessageResponse{}
		msgRes.Message = ""
		return msgRes, nil
	}

	var refresh MessageResponse
	err = json.NewDecoder(res.Body).Decode(&refresh)
	if err != nil {
		fmt.Printf("error decoding refresh request: %s\n", err)
	}
	return &refresh, err
}

func (httpClient *HttpGameClient) GetLobby() ([]PlayerStatus, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/lobby"

	req, _ := http.NewRequest("GET", requestURL, nil)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error getting lobby request: %s\n", err)
	}
	defer res.Body.Close()

	lobby := []PlayerStatus{}
	err = json.NewDecoder(res.Body).Decode(&lobby)
	if err != nil {
		fmt.Printf("error decoding lobby request: %s\n", err)
	}
	return lobby, err
}

func (httpClient *HttpGameClient) GetTop10Players() (Top10Players, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/stats"

	req, _ := http.NewRequest("GET", requestURL, nil)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error getting Top 10 players request: %s\n", err)
	}
	defer res.Body.Close()

	var tp Top10Players
	err = json.NewDecoder(res.Body).Decode(&tp)
	if err != nil {
		fmt.Printf("error decoding Top 10 playersrequest: %s\n", err)
	}
	return tp, err
}

func (httpClient *HttpGameClient) GetPlayerStats(nick string) (*PlayerStats, error) {
	requestURL := fmt.Sprintf("https://go-pjatk-server.fly.dev/api/stats/%v", nick)

	req, _ := http.NewRequest("GET", requestURL, nil)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error getting players stats: %s\n", err)
	}
	defer res.Body.Close()

	var ps PlayerStats
	err = json.NewDecoder(res.Body).Decode(&ps)
	if err != nil {
		fmt.Printf("error decoding player stats: %s\n", err)
	}
	return &ps, err
}

func (httpClient *HttpGameClient) AbandonGame() (*MessageResponse, error) {
	requestURL := "https://go-pjatk-server.fly.dev/api/game/abandon"

	req, _ := http.NewRequest(http.MethodDelete, requestURL, nil)
	req.Header.Set("X-Auth-Token", httpClient.XAuthToken)
	res, err := httpClient.makeRequest(req)

	if err != nil {
		fmt.Printf("error abandoning: %s\n", err)
	}
	defer res.Body.Close()

	var message MessageResponse

	if res.StatusCode == 200 {
		msgRes := &MessageResponse{}
		msgRes.Message = ""
		return msgRes, nil
	}

	err = json.NewDecoder(res.Body).Decode(&message)
	if err != nil {
		fmt.Printf("error decoding abandoning: %s\n", err)
	}
	return &message, err
}

func handleResponseCode(statusCode int, req *http.Request) string {
	var httpResposneError string
	switch statusCode {
	case 401:
		httpResposneError = "Unauthorized: 401"
	case 400:
		httpResposneError = "Bad Request: 400"
	case 403:
		httpResposneError = "Forbidden: 403"
	case 429:
		httpResposneError = fmt.Sprintf("Too Many Requests: 429 URL: %s", req.URL)
	case 503:
		httpResposneError = "Service Unavailable: 503"
	default:
		httpResposneError = fmt.Sprintf("Unhandled status code: %d\n", statusCode)
	}
	return httpResposneError
}
