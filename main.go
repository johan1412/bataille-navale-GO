package main

import (
	"fmt"
	"log"
	"net/http"
)

type square struct {
	x   int
	y   int
	hit bool
}

type ship struct {
	squares []square
}

func (ship *ship) isSunk() bool {
	i := 0
	for i < len(ship.squares) {
		if !ship.squares[i].hit {
			return false
		}
		i++
	}
	return true
}

func main() {
	fmt.Println("=================\n BATLESHIP START \n=================")

	mux := http.NewServeMux()

	mux.HandleFunc("/board", BoardHandler)
	mux.HandleFunc("/boats", BoatsHandler)
	mux.HandleFunc("/hit", HitHandler)

	WrappedMux := RunSomeCode(mux)
	log.Fatal(http.ListenAndServe(":4567", WrappedMux))
}
