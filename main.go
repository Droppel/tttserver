package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"tttserver/game"
)

const (
	invalidSyntaxError = "INVALID_SYNTAX"
	lobbyDoesntExist   = "LOBBY_DOES_NOT_EXIST"
)

const (
	connHost = "0.0.0.0"
	connPort = "3333"
	connType = "tcp"
)

var (
	lobbies map[string]*game.Lobby
)

func main() {
	lobbies = make(map[string]*game.Lobby)

	// Listen for incoming connections.
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + connHost + ":" + connPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(game.InitClient(conn))
	}
}

// Handles incoming requests.
func handleRequest(client *game.Client) {
	for {
		bufReader := bufio.NewReader(client)
		// Read tokens delimited by newline
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		command := strings.Split(string(bytes[:len(bytes)-1]), " ")
		switch command[0] {
		case "create":
			if len(command) != 3 {
				client.Write([]byte("ERROR\n"))
				continue
			}
			size, err := strconv.Atoi(command[1])
			if err != nil {
				client.Write([]byte("ERROR\n"))
				continue
			}
			sizeWin, err := strconv.Atoi(command[2])
			if err != nil {
				client.Write([]byte("ERROR\n"))
				continue
			}
			playerCount := 2
			// playerCount, err := strconv.Atoi(command[3])
			// if err != nil {
			// 	client.Write([]byte("ERROR\n"))
			// 	continue
			// }

			newLobby := game.CreateLobby(size, sizeWin, playerCount)
			newLobby.SetClient(0, client)
			lobbies[newLobby.GetId()] = newLobby
			client.WriteMsg(newLobby.GetId())
		case "join":
			if len(command) != 2 {
				client.WriteMsg("ERROR " + invalidSyntaxError)
				continue
			}
			lobby, exists := lobbies[command[1]]
			if !exists {
				client.WriteMsg("ERROR " + lobbyDoesntExist)
				continue
			}

			success := lobby.SetClient(1, client)
			if !success {
				client.WriteMsg("ERROR")
				continue
			}
			client.WriteMsg("SUCCESS")
			lobby.StartGame()
			lobby.SendBoardState()
		case "move":
			if len(command) != 3 {
				client.WriteMsg("ERROR")
				continue
			}
			x, err := strconv.Atoi(command[1])
			if err != nil {
				client.WriteMsg("ERROR")
				continue
			}
			y, err := strconv.Atoi(command[2])
			if err != nil {
				client.WriteMsg("ERROR")
				continue
			}

			lobby := lobbies[client.GetLobby()]
			move := game.InitPoint(x, y)
			success := lobby.MakeMove(client, move)
			if success {
				client.WriteMsg("SUCCESS")
			} else {
				client.WriteMsg("ERROR")
				continue
			}
			done, win := lobby.CheckDone(move)
			if done {
				lobby.EndGame(win)
				delete(lobbies, lobby.GetId())
			} else {
				lobby.SendBoardState()
			}
		default:
			client.WriteMsg("ERROR")
		}

	}
}
