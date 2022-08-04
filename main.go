package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Store struct {
	Log *logrus.Logger
}

func NewStore() (*Store, error) {

	return &Store{
		Log: logrus.New(),
	}, nil
}

func genLink() (string, error) {
	link, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return link.String(), nil
}

func createFile(name string) (string, error) {
	filename := name + ".go"
	cmd := exec.Command("touch", filename)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *Store) prepareFile(data []byte) (string, error) {
	link, err := genLink()
	if err != nil {
		s.Log.Error(err)
		return "", err
	}

	filename, err := createFile(link)
	if err != nil {
		s.Log.Error(err)
		return "", err
	}

	err = ioutil.WriteFile(filename, data, 0700)
	if err != nil {
		s.Log.Error(err)
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", "go_binary_files", filename)
	err = cmd.Run()
	if err != nil {
		s.Log.Error(err)
		return "", err
	}

	cmd = exec.Command("mv", filename, "go_files")
	err = cmd.Run()
	if err != nil {
		s.Log.Error(err)
		return "", err
	}

	return link, nil
}

func (s *Store) Run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, code, err := conn.ReadMessage()
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	link, err := s.prepareFile(code)
	if err != nil {
		s.Log.Error(err)
		err = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("./go_binary_files/" + link)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = cmd.Start()
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := bufio.NewReader(pipe)

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			conn.Close()
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, line)
		if err != nil {
			s.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *Store) Share(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, code, err := conn.ReadMessage()
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	link, err := s.prepareFile(code)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(link))
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Store) PlayGo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, "index.html")
}

func (s *Store) GetCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	link := ps.ByName("link")
	data, err := ioutil.ReadFile("./go_files/" + link + ".go")
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(data)
	http.ServeFile(w, r, "index.html")
}

func main() {
	store, err := NewStore()
	if err != nil {
		store.Log.Fatal(err)
	}

	router := httprouter.New()

	router.GET("/", store.PlayGo)
	router.GET("/play", store.Run)
	router.GET("/share", store.Share)
	router.GET("/p/:link", store.GetCode)

	err = http.ListenAndServe(":3000", router)
	store.Log.Fatal(err)
}
