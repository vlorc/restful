package main

import (
	"errors"
	"github.com/vlorc/restful/pkg/engine"
	"log"
	"net/http"
)

type EchoRequest struct {
	Data string `json:"data"`
}

type EchoResponse struct {
	Data interface{} `json:"data"`
}

func main() {
	engine.Init()

	g := engine.NewRouter()

	g.Any("/echo", func(req *EchoRequest) (*EchoResponse, error) {
		log.Println("body: ", req.Data)

		if len(req.Data) == 0 {
			return nil, errors.New("Illegal parameter")
		}
		return &EchoResponse{Data: req.Data}, nil
	})

	http.ListenAndServe(":1234", g)
}
