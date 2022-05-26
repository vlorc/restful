package main

import (
	"errors"
	"github.com/vlorc/restful/pkg/engine"
	_ "github.com/vlorc/restful/pkg/engine/chi"
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

	g := engine.Default()

	g.Any("/echo", func(req *EchoRequest) (*EchoResponse, error) {
		log.Println("body: ", req.Data)

		if len(req.Data) == 0 {
			return nil, errors.New("Illegal parameter")
		}
		return &EchoResponse{Data: req.Data}, nil
	})

	http.ListenAndServe(":1234", g)
}
