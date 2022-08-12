package code

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/alekstet/my_playground/files"
	"github.com/gorilla/websocket"
)

type Coder interface {
	GetCode(link string) ([]byte, error)
	ShareCode(conn *websocket.Conn) error
	Run(conn *websocket.Conn) (io.ReadCloser, error)
}

func (s CodeStore) GetCode(link string) ([]byte, error) {
	data, err := os.ReadFile("./" + s.Config.FilesFolder + "/" + link + ".go")
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(string(data))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s CodeStore) Run(conn *websocket.Conn) (io.ReadCloser, error) {
	_, code, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	link, err := files.PrepareFile(code, s.Config)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return nil, err
	}

	cmd := exec.Command("./" + s.Config.BinariesFolder + "/" + link)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return pipe, nil
}

func (s CodeStore) ShareCode(conn *websocket.Conn) error {
	_, code, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	link, err := files.PrepareFile(code, s.Config)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(link))
	if err != nil {
		return err
	}

	return nil
}
