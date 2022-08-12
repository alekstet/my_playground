package main

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	linkCode    string
	testProgram = `
	// You can edit this code!
	// Click here and start typing.
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, PlayGO")
	}`
)

type Code string

func TestRun(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/run", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(testProgram)); err != nil {
		t.Fatalf("%v", err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	assertions := assert.New(t)
	assertions.Equal("Hello, PlayGO", string(p))
}

func TestShare(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/share", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(testProgram)); err != nil {
		t.Fatalf("%v", err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	linkCode = string(p)
}

func TestGetCode(t *testing.T) {
	client := resty.New()
	r, err := client.R().Get("http://localhost:8080/pp/" + linkCode)
	if err != nil {
		t.Error(err)
	}

	var code Code
	json.Unmarshal(r.Body(), &code)
	assertions := assert.New(t)
	assertions.EqualValues(code, testProgram)
}

func TestMain(t *testing.T) {
	cnf := "config/config.yml"
	go Run(&cnf)

	TestShare(t)
	TestGetCode(t)
	TestRun(t)
}
