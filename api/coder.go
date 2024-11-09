package api

import (
	"bufio"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func (s Store) getCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	link := ps.ByName("link")

	resp, err := s.Coder.GetCode(link)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (s Store) pipe(pipe io.ReadCloser, conn *websocket.Conn, w http.ResponseWriter) {
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
func (s Store) playGo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, s.Config.HTMLName)
}

func (s Store) run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pipe, err := s.Coder.Run(conn)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.pipe(pipe, conn, w)
}

func (s Store) share(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Coder.ShareCode(conn)
	if err != nil {
		s.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
