package goServer

import (
	"fmt"
	"log"
	"net/http"
	"golang.org/x/net/websocket"
	"encoding/json"
	"os"
)

type Message struct {
	ID     int
	Md5Sum [16]byte
	Body   []byte
}

func wsRepeat(ws *websocket.Conn) {
	var err error
	var reply []byte
	var msg Message

	if err = websocket.Message.Receive(ws, &reply); err != nil {
		fmt.Println(err)
	}
	if err = websocket.Message.Send(ws, "OK"); err != nil {
		fmt.Println("Cannot send ws message")
		log.Fatal(err)
	}

	fmt.Println(reply)

	err = json.Unmarshal(reply, &msg)
	check(err)

	f, err := os.OpenFile("output.txt", os.O_APPEND | os.O_WRONLY, 0600)
	check(err)

	_, err = f.Write(msg.Body)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	f, err := os.Create("output.txt")
	check(err)
	defer f.Close()
	http.Handle("/", websocket.Handler(wsRepeat))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}

}


