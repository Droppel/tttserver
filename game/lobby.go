package game

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	net.Conn
	lobby  string
	number int
}

func InitClient(conn net.Conn) *Client {
	return &Client{conn, "", 0}
}

func (c *Client) GetLobby() string {
	return c.lobby
}

func (c *Client) setLobby(id string) {
	c.lobby = id
}

func (c *Client) WriteMsg(msg string) (int, error) {
	return c.Write([]byte(msg + "\n"))
}

type Lobby struct {
	id          string
	g           *Game
	playerCount int
	clients     []*Client
}

func CreateLobby(size, sizeWin, playerCount int) *Lobby {
	id := uuid.NewString()
	return &Lobby{id: id, g: InitGame(size, sizeWin), playerCount: playerCount, clients: make([]*Client, 2)}
}

func (l *Lobby) BroadcastMsg(msg string) {
	for _, client := range l.clients {
		client.WriteMsg(msg)
	}
}

func (l *Lobby) SendBoardState() error {
	var client *Client
	if l.GetActivePlayerNum() == l.clients[0].number {
		client = l.clients[0]
	} else {
		client = l.clients[1]
	}

	_, err := client.WriteMsg(l.g.StateToString())
	if err != nil {
		return err
	}
	return nil
}

func (l *Lobby) MakeMove(c *Client, move Point) bool {
	if c.number != l.GetActivePlayerNum() {
		return false
	}
	return l.g.MakeMove(move)
}

func (l *Lobby) CheckDone(move Point) (bool, int) {
	return l.g.CheckDone(move)
}

func (l *Lobby) GetActivePlayerNum() int {
	return l.g.GetActivePlayerNum()
}

func (l *Lobby) GetId() string {
	return l.id
}

func (l *Lobby) SetClient(num int, c *Client) bool {
	if l.clients[num] != nil {
		return false
	}
	l.clients[num] = c
	c.setLobby(l.GetId())
	return true
}

func (l *Lobby) StartGame() {
	if rand.Intn(2) == 0 {
		l.clients[0].number = 1
		l.clients[1].number = 2
	} else {
		l.clients[0].number = 2
		l.clients[1].number = 1
	}

	msg := fmt.Sprintf("START %d", l.g.size)
	for _, client := range l.clients {
		client.WriteMsg(fmt.Sprintf("%s %d", msg, client.number))
	}
}

func (l *Lobby) EndGame(win int) {
	l.BroadcastMsg(fmt.Sprintf("FINISH %d", win))
	for _, client := range l.clients {
		client.Close()
	}
}
