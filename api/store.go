package api

import (
	"github.com/alekstet/my_playground/code"
	"github.com/alekstet/my_playground/config"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Store struct {
	Log      *logrus.Logger
	Config   *config.PlayConfig
	Upgrader websocket.Upgrader
	Coder    code.Coder
}

func NewStore(cfg config.PlayConfig, code code.CodeStore) (*Store, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	store := &Store{
		Log:      logrus.New(),
		Config:   &cfg,
		Upgrader: upgrader,
	}
	store.Coder = code

	return store, nil
}
