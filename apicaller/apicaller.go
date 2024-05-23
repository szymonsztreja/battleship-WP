package apicaller

import (
	client "battleship-WP/temp"
	"fmt"
	"net/http"
)

func Call() {
	client := &client.HttpGameClient{
		Client: &http.Client{},
	}

	for i := 0; i < 5; i++ {
		// statusCode, err := client.InitGame()
		// fmt.Println("InitGame:", statusCode, err)

		// _, statusCode, err = client.Board()
		// fmt.Println("Board:", statusCode, err)

		// _, statusCode, err = client.Status()
		// fmt.Println("Status:", statusCode, err)

		// _, statusCode, err = client.Fire("A2")
		// fmt.Println("Fire:", statusCode, err)

		// _, statusCode, err = client.GetPlayersDescription()
		// fmt.Println("GetPlayersDescription:", statusCode, err)

		lobby, err := client.GetLobby()
		fmt.Println("GetLobby:", lobby, err)
	}
}
