package api

import "github.com/julienschmidt/httprouter"

func (s *Store) Register(router *httprouter.Router) {
	//HTTP handlers
	router.GET("/", s.playGo)
	router.GET("/p/:link", s.playGo)
	router.GET("/pp/:link", s.getCode)

	//Websocket handlers
	router.GET("/run", s.run)
	router.GET("/share", s.share)
}
