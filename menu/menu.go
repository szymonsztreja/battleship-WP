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
	"sync"
	"time"
)

type Menu struct {
	// httpClient := &client.HttpGameClient{
	// 	Client: &http.Client{},
	// }
	httpClient *client.HttpGameClient
	game       *game.Game
	player     *player
}

type player struct {
	nick string
	desc string
}

func (m *Menu) Start() {
	m.httpClient = &client.HttpGameClient{
		Client: &http.Client{},
	}
	m.player = &player{}

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
			fmt.Println("You chose: Set name and description")
			m.setNickAndDesc()
		case 2:
			fmt.Println("You chose: Play a game")
			m.playModes()
		case 3:
			fmt.Println("You chose: View Top 10 Players statistics")
			m.getTop10Players()
		case 4:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please choose a number between 1 and 4.")
		}
	}

}

func (m *Menu) playModes() {

	for {
		fmt.Println("Set your game mode!")
		fmt.Println("1. Play with bot")
		fmt.Println("2. Play with a player in lobby")
		fmt.Println("3. Get challenged by a player")
		fmt.Println("4. Go back")

		choice := playerInput("")
		option, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Print(err)
		}
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
			fmt.Println("Returning")
			return
		}
	}
}

func (m *Menu) playWithBot() {
	m.createGameInstance()
	m.game.Wpbot = true
	m.gameSetup()
	fmt.Println("before watiing for game")
	m.waitForGame()
	fmt.Println("starting game")
	m.game.Run()
}

func (m *Menu) playWithPlayer() {
	// List players in a lobby
	wp := m.waitingPlayers()
	if len(wp) == 0 {
		fmt.Println("Empty lobby")
		return
	}

	for _, player := range wp {
		fmt.Println(player.Nick, player.GameStatus)
	}

	tn := m.handlePlayerChallenge(wp)
	m.createGameInstance()
	m.game.TargetNick = tn
	m.gameSetup()
	m.waitForGame()
	m.game.Run()
}

func (m *Menu) setNickAndDesc() {
	nick := playerInput("Type in your nick")
	desc := playerInput("Type in your description")

	m.player.nick = nick
	m.player.desc = desc
}

func (m *Menu) getChellengedByPlayer() {
	m.createGameInstance()
	m.game.Wpbot = false
	m.gameSetup()
	m.waitingForChallenge()
	m.game.Run()
}

func (m *Menu) waitingForChallenge() bool {
	// Create an unbuffered channel to receive the result
	done := make(chan bool)
	ticker10Sec := time.NewTicker(10 * time.Second)
	ticker1Sec := time.NewTicker(1 * time.Second)
	var wg sync.WaitGroup

	// Increment WaitGroup counter
	wg.Add(2)

	// Status goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				ticker1Sec.Stop()
				//return // Exit if signaled

			case <-ticker1Sec.C:
				fmt.Println("getting status goroutine")
				status, err := m.httpClient.Status()
				if err != nil {
					fmt.Println("Error waiting for challenge:", err)
				}
				if status.GameStatus == "game_in_progress" {
					// Send true when the game is in progress
					done <- true
					fmt.Println("Game is in progress")
					// case <-time.After(time.Second * 1):
					// 	fmt.Println("Timeout sending true to done channel")
					// return
				}
			}
		}
	}()

	// Refreshing goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				ticker10Sec.Stop()
				fmt.Println("exiting refreshing")
				// return // Exit if signaled
			case <-ticker10Sec.C:
				fmt.Println("Refreshing session!")
				res, err := m.httpClient.RefreshSession()
				if err != nil {
					fmt.Printf("Error refreshing session message: %v\n", res.Message)
				}
			}
		}
	}()

	// Ensure cleanup even if the function exits prematurely
	defer func() {
		close(done) // Signal all goroutines to exit
		wg.Wait()   // Wait for all goroutines to finish
	}()

	// Wait for a value on the channel
	return <-done
}

func (m *Menu) getTop10Players() {
	top10, err := m.httpClient.GetTop10Players()
	if err != nil {
		fmt.Print("Error getting top players statistic")
	}
	for _, player := range top10.Stats {
		fmt.Println("--------------------------------------")
		fmt.Println("Games | Nick \t| Points | Rank | Wins")
		fmt.Print(player.Games, "\t")
		fmt.Print(player.Nick, "\t")
		fmt.Print(player.Points, "\t")
		fmt.Print(player.Rank, "\t")
		fmt.Print(player.Wins, "\t\n")
	}
}

func (m *Menu) createGameInstance() {
	m.game = &game.Game{
		HttpGameC: m.httpClient,
	}
}

func (m *Menu) gameSetup() {
	gameData := client.GameData{
		Coords:     []string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"},
		Desc:       m.player.desc,
		Nick:       m.player.nick,
		TargetNick: m.game.TargetNick,
		Wpbot:      m.game.Wpbot,
	}
	m.httpClient.InitGame(gameData)
}

func (m *Menu) waitForGame() {
	for {
		status, err := m.httpClient.Status()
		if err != nil {
			fmt.Printf("error getting game status : %s\n", err)
		}

		if status.GameStatus == "game_in_progress" {
			fmt.Println(status.GameStatus)
			return
		}

		time.Sleep(1 * time.Second)
	}
}

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
		tn = challengePlayerToDuel()
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
