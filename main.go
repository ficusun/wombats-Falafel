package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"fmt"
	"log"
	"net/http"
)

var Hub ClientHub = ClientHub{
	Connections: make(map[int]*Client, 0),
}

var MessageForUsers = make(chan []byte, 1000)

var MessageFromUsers = make(chan []byte, 1000)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	for {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ServeWs, upgrader.Upgrade: ", err)
			return
		}

		fmt.Println("new client: ", conn.RemoteAddr())
		// add conn to

		res := Hub.AddClient(InitClient(conn))
		if res {
			fmt.Println("successful")
		} else {
			fmt.Println("failed to add client")
		}

	}
}

func Start() {

	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		ServeWs(w, r)
	})

	var addr = flag.String("addr", "192.168.0.52:8081", "http service address")

	err := http.ListenAndServe(*addr, nil) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}


func Router (MessageFromUsers chan []byte) {
	for message := range MessageFromUsers {

		letter := Letter{}

		err := json.Unmarshal(message, &letter)
		if err != nil {
			log.Println(err)
		}

		messageForUser, err := json.Marshal(letter.Letter)
		if err != nil {
			log.Println(err)
		}

		for id, _ := range Hub.Connections {
			if id != letter.ClientID {
				Hub.Connections[id].Send<- messageForUser //[]byte(letter)
			}
		}
	}
}

func main() {
	fmt.Println("start app")
	go Router(MessageFromUsers)
	Start()
}
