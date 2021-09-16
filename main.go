package main

import (
	"math/rand"
	"time"
	"fmt"
	"log"
	"net/http"
	"os"
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
	size := 10
	var takenSquares [10][10]bool
	var ships [3]ship
	i := 0
	for i < len(ships) {

		rand.Seed(time.Now().UnixNano())
		direction := rand.Intn(2) // 0 = horizontal & 1 = vertical

		var xmax int
		var ymax int

		if direction == 0 {
			xmax = size - i
			ymax = size
		} else {
			xmax = size
			ymax = size - i
		}

		x := rand.Intn(xmax)
		y := rand.Intn(ymax)

		if canPutShip(direction, x, y, i+1, takenSquares) {
			j := 0
			for j <= i {
				if direction == 0 {
					var square square
					square.x = x + j
					square.y = y
					takenSquares[x+j][y] = true
					ships[i].squares = append(ships[i].squares, square)

				} else {
					var square square
					square.x = x
					square.y = y + j
					takenSquares[x][y+j] = true
					ships[i].squares = append(ships[i].squares, square)
				}

				j++
			}
		} else {
			//fmt.Print("continue\n")
			continue
		}
		i++

	}

	port := os.Args[1]
	fmt.Println("Port = ", port)

	listenPort := ":" + port

	//mux := http.NewServeMux()

	http.HandleFunc("/board", BoardHandler)
	http.HandleFunc("/boats", BoatsHandler)
	http.HandleFunc("/hit", HitHandler)

	//WrappedMux := RunSomeCode(mux)
	log.Fatal(http.ListenAndServe(listenPort, nil))




	//boucle commandes : connect, attack

}

func canPutShip(direction int, x int, y int, size int, takenSquares [10][10]bool) bool {

	if direction == 0 { //h

		for i := 0; i < size; i++ {
			if takenSquares[x+i][y] {
				return false
			}
		}

	} else { //v

		for i := 0; i < size; i++ {
			if takenSquares[x][y+i] {
				return false
			}
		}

	}

	return true
}

