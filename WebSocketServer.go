package main

import (
	"fmt"
	"log"
	"net/http"
	"golang.org/x/net/websocket"
	"encoding/json"
	"os"
	"crypto/md5"
	"bytes"
)

type Message struct {
	ID     int
	Md5Sum [16]byte
	Body   []byte
}

var f *os.File

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
	md5Sum := md5.Sum(msg.Body)
	if (bytes.Equal(md5Sum[:], msg.Md5Sum[:])) {
		writeFile(msg.Body)
	}
}

func writeFile(message []byte) {
	_, err := f.Write(message)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func openFile(fileName string) {
	var err error
	f, err = os.OpenFile(fileName, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	check(err)
}

func main() {
	openFile("output.txt")
	defer f.Close()
	http.Handle("/", websocket.Handler(wsRepeat))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

