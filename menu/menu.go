package menu

import (
	"battleship-WP/client"
	"battleship-WP/game"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Menu struct {
	// httpClient := &client.HttpGameClient{
	// 	Client: &http.Client{},
	// }
	httpClient *client.HttpGameClient
	game       *game.Game
}

func (m *Menu) Start() {
	m.httpClient = &client.HttpGameClient{
		Client: &http.Client{},
	}
	m.game = &game.Game{}

	// Set your nick and description
	// nick := playerInput("nick")
	// desc := playerInput("description")
	// println(nick, desc)

	// game := game.Game{
	// 	PlayerNick:        nick,
	// 	PlayerDescription: desc,
	// 	TargetNick:        "",
	// 	Wpbot:             true,
	// }
	// game.Run()
	for {
		fmt.Println("Welcome to the Command Line Menu!")
		fmt.Println("1. Set name and description")
		fmt.Println("2. Play a game")
		fmt.Println("3. View Top 10 Players statistics")
		fmt.Println("4. Exit")

		choice := playerInput("")
		option, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Print(err)
		}
		switch option {
		case 1:
			fmt.Println("You chose Option 1")
			m.setNickAndDesc()
		case 2:
			fmt.Println("You chose Option 2")
			m.play()
		case 3:
			fmt.Println("You chose Option 3")
			// Add your logic for Option 3 here
		case 4:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please choose a number between 1 and 4.")
		}
	}

}

func (m *Menu) play() {

	fmt.Println("Set your game mode!")
	fmt.Println("1. Play with bot")
	fmt.Println("2. Play with a player in lobby")
	fmt.Println("3. Get challenged by a player")

	choice := playerInput("")
	option, err := strconv.Atoi(choice)
	if err != nil {
		fmt.Print(err)
	}

	for {
		switch option {
		case 1:
			fmt.Println("Play with bot")
			m.playWithBot()
		case 2:
			fmt.Println("Play with a player in lobby")
			m.playWithPlayer()
		case 3:
			fmt.Println("Get challenged by a player")
			m.getChellengedByPlayer()
		case 4:
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}

func (m *Menu) playWithBot() {
	m.game.Wpbot = true
	m.game.Run()
}

func (m *Menu) playWithPlayer() {
	// List players in a lobby
	wp := m.waitingPlayers()
	if len(wp) != 0 {
		for _, player := range wp {
			fmt.Println(player.Nick, player.GameStatus)
		}
	} else {
		fmt.Println("Empty lobby")
	}

	tn := m.handlePlayerChallenge(wp)

	m.game.TargetNick = tn
	m.game.Run()

}

func (m *Menu) getChellengedByPlayer() {
	m.game.Wpbot = false
	m.waitingForChallenge()
	m.game.Run()

}

func (m *Menu) setNickAndDesc() {
	nick := playerInput("Type in your nick")
	desc := playerInput("Type in your description")

	m.game.PlayerNick = nick
	m.game.PlayerDescription = desc
}

func (m *Menu) waitingForChallenge() bool {
	// Create a channel to receive the result
	statusChan := make(chan bool)

	go func() {
		for {
			m.httpClient.RefreshSession()
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			status, err := m.httpClient.Status()
			if err != nil {
				fmt.Println("Error waiting for challenge:", err)
				statusChan <- false // Send false in case of error
				return
			}
			if status.GameStatus == "game_in_progress" {
				statusChan <- true // Send true when the game is in progress
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()

	// Wait for a value on the channel
	return <-statusChan
}

// func (m *Menu) getTop10Players(){
// 	top10, err := m.httpClient.GetTop10Players()
// 	if err != nil {
// 		fmt.Print("Error getting top players statistic")
// 	}
// }

func getYesOrNo(prompt string) bool {
	// reader := bufio.NewReader(os.Stdin)
	var output bool
	fmt.Println(prompt)
	for {
		yesNo := playerInput("answer")

		// Trim any leading or trailing whitespace
		input := strings.TrimSpace(yesNo)

		// Validate input
		if input != "Y" && input != "N" {
			fmt.Println("Please enter Y or N")
			continue
		}

		if input == "Y" {
			output = true
		}

		if input == "N" {
			output = false
		}

		return output
	}
}

func (m *Menu) handlePlayerChallenge(wp []client.PlayerStatus) string {
	var tn string

	for {
		tn := challengePlayerToDuel()
		found := false

		for _, player := range wp {
			if player.Nick == tn {
				found = true
				break
			} else {
				fmt.Printf("No player named:%v in a lobby\n", tn)
			}
		}

		if found {
			break
		}
	}
	return tn
}

func playerInput(thingToType string) string {
	fmt.Println(thingToType)
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

func challengePlayerToDuel() string {
	// challengeText := fmt.Sprint("enemies nick to challenge")
	target_nick := playerInput("enemy nick to challenge: ")
	return target_nick
}
