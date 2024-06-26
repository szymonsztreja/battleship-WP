package menu

import (
	"battleship-WP/client"
	"battleship-WP/game"
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
	// placeShips "battleship-WP/"
)

var defaultShips = []string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"}

type Menu struct {
	// httpClient := &client.HttpGameClient{
	// 	Client: &http.Client{},
	// }
	httpClient *client.HttpGameClient
	game       *game.Game
	player     *player
	ships      []string
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
	m.ships = []string{}

	for {
		// ClearScreen()
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
			m.playOrPlace()
		case 3:
			fmt.Println("You chose: View Top 10 Players statistics")
			m.getTop10Players()
		case 4:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please choose a number. According to menu options")
		}
	}

}

func (m *Menu) playOrPlace() {
	for {
		// ClearScreen()
		fmt.Println("Set your ships or select a game mode!")
		fmt.Println("1. Game modes")
		fmt.Println("2. Set ships")
		fmt.Println("3. Go back")

		choice := playerInput("")
		option, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Print(err)
		}
		switch option {
		case 1:
			fmt.Println("Games modes")
			m.playModes()
		case 2:
			m.shipsMenu()
		case 3:
			fmt.Println("Returning")
			return
		default:
			fmt.Println("Invalid choice. Please choose a number. According to menu options")
		}
	}
}

func (m *Menu) shipsMenu() {
	for {
		// ClearScreen()
		fmt.Println("Set your ships!")
		fmt.Println("1. Play with default ship placement")
		fmt.Println("2. Set your own ships")
		fmt.Println("3. Check ship placement")
		fmt.Println("4. Go back to play or place")

		choice := playerInput("")
		option, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Print(err)
		}
		switch option {
		case 1:
			fmt.Println("Default ships")
			m.ships = defaultShips
		case 2:
			fmt.Println("Setting your own ships")
			m.ships = game.PlaceShips()
		case 3:
			if len(m.ships) < 20 {
				fmt.Println("Ships not set or set incorectly")
			} else {
				fmt.Println("Ship set corectly. Ready to rumble")
			}
		case 4:
			fmt.Println("Returning")
			return
		default:
			fmt.Println("Invalid choice. Please choose a number. According to menu options")

		}
	}
}

func (m *Menu) playModes() {

	for {
		// ClearScreen()
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
		default:
			fmt.Println("Invalid choice. Please choose a number. According to menu options")
		}
	}
}

func (m *Menu) playWithBot() {
	m.createGameInstance()
	m.game.Wpbot = true
	err := m.gameSetup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m.ships)
	fmt.Println(m.httpClient.XAuthToken)
	m.waitForGame()
	fmt.Println("starting game")
	m.game.Run()
	m.resetGameInstance()
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
	err := m.gameSetup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.waitForGame()
	m.game.Run()
	m.resetGameInstance()

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
	err := m.gameSetup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.waitingForChallenge()
	m.game.Run()
	m.resetGameInstance()
}

func (m *Menu) waitingForChallenge() bool {
	// Create an unbuffered channel to receive the result
	gameInProgressCh := make(chan bool)
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
			case <-gameInProgressCh:
				ticker1Sec.Stop()
				return

			case <-ticker1Sec.C:
				fmt.Println("getting status goroutine")
				status, err := m.httpClient.Status()
				if err != nil {
					fmt.Println("Error waiting for challenge:", err)
				}
				if status.GameStatus == "game_in_progress" {
					// Send true when the game is in progress
					gameInProgressCh <- true
					fmt.Println("Game is in progress")
				}
			}
		}
	}()

	// Refreshing goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-gameInProgressCh:
				ticker10Sec.Stop()
				fmt.Println("exiting refreshing")
				return
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
		close(gameInProgressCh) // Signal all goroutines to exit
		wg.Wait()               // Wait for all goroutines to finish
	}()

	// Wait for a value on the channel
	return <-gameInProgressCh
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

func (m *Menu) resetGameInstance() {
	m.game = nil
}

func (m *Menu) gameSetup() error {
	var err error

	gameData := client.GameData{
		Coords:     m.ships,
		Desc:       m.player.desc,
		Nick:       m.player.nick,
		TargetNick: m.game.TargetNick,
		Wpbot:      m.game.Wpbot,
	}
	mes := m.httpClient.InitGame(gameData)

	if mes != "" {
		err = errors.New(mes)
		return err
	}
	return err
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
	target_nick := playerInput("enemy nick to challenge: ")
	return target_nick
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
