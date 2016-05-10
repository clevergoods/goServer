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
var c chan Message
var counter int
var m map[int][]byte

func wsRepeat(ws *websocket.Conn) {
	var err error
	var reply []byte

	if err = websocket.Message.Receive(ws, &reply); err != nil {
		fmt.Println(err)
	}

	if err = websocket.Message.Send(ws, "OK"); err != nil {
		fmt.Println("Cannot send ws message")
		log.Fatal(err)
	}

	go checkMd5Sum(reply)
	go writeFile()
}

func writeFile() {
	for s := range (c) {
		if counter == s.ID {
			msg := s.Body
			_, err := f.Write(msg)
			check(err)
			fmt.Println("from chan", counter)
			fok := true

			counter++
			for ; fok; {
				msg, fok = m[counter]
				if fok {
					_, err := f.Write(msg)
					check(err)
					fmt.Println("write from m", counter)
					delete(m, counter)
					counter++
				}
			}

		}else {
			m[s.ID] = s.Body
			fok := true
			var msg []byte
			for ; fok; {
				msg, fok = m[counter]
				if fok {
					_, err := f.Write(msg)
					check(err)
					fmt.Println("from m", counter)
					delete(m, counter)
					counter++
				}
			}
		}
	}
}

func checkMd5Sum(reply []byte) {
	var msg Message
	err := json.Unmarshal(reply, &msg)
	check(err)
	md5Sum := md5.Sum(msg.Body)
	if (bytes.Equal(md5Sum[:], msg.Md5Sum[:])) {
		c <- msg
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func openFile(fileName string) {
	var err error
	counter = 0
	f, err = os.OpenFile(fileName, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	check(err)
}

func main() {
	c = make(chan Message, 100)
	m = make(map[int][]byte)
	counter = 0
	openFile("output.txt")
	defer f.Close()
	http.Handle("/", websocket.Handler(wsRepeat))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

