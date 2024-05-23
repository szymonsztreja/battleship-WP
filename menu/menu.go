package menu

import (
	"battleship-WP/client"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Menu struct {
	// httpClient := &client.HttpGameClient{
	// 	Client: &http.Client{},
	// }
	httpClient *client.HttpGameClient
}

func (m *Menu) Start() {
	m.httpClient = &client.HttpGameClient{
		Client: &http.Client{},
	}

	nick := playerInput("nick")
	desc := playerInput("description")
	println(nick, desc)

	wp := m.waitingPlayers()

	for _, player := range wp {
		fmt.Println(player)
	}
	// game := game.Game{
	// 	PlayerNick:        nick,
	// 	PlayerDescription: desc,
	// }
	// game.Run()
}

func playerInput(thingToType string) string {
	fmt.Printf("Type in your %v:\n", thingToType)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return scanner.Text()
}

func (m *Menu) waitingPlayers() []client.PlayerStatus {
	wp, err := m.httpClient.GetLobby()
	if err != nil {
		fmt.Printf("error getting waiting players : %s\n", err)
	}
	return wp
}
