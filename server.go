package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func BoardHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}

		stringBoard := ""

		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				stringBoard += strconv.Itoa(board[i][j]) + " | "
			}
			stringBoard += "\n"
		}

		fmt.Fprintf(w, stringBoard)
	}
}

func BoatsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		remainingShips := 0
		message := ""

		for i := 0; i < len(ships); i++ {
			if !ships[i].isSunk() {
				remainingShips++
			}
		}

		if remainingShips == 0 {
			message = "\nPlayer lost"
		}

		fmt.Fprintf(w, "Remaining ships : "+strconv.Itoa(remainingShips)+message)
	}
}

func HitHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodPost:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}

		x, _ := strconv.Atoi(req.PostForm["x"][0])
		y, _ := strconv.Atoi(req.PostForm["y"][0])
		x--
		y--
		message := ""

		if !takenSquares[x][y] {
			board[x][y] = 1
			message = "You didn't touch any ship"

		} else {
			board[x][y] = 2
		out:
			for i := 0; i < len(ships); i++ {
				message = "You touched a ship !"
				for j := 0; j < len(ships[i].squares); j++ {
					if ships[i].squares[j].x == x && ships[i].squares[j].y == y {
						ships[i].squares[j].hit = true
						if ships[i].isSunk() {
							message = "You sunk a ship !"
						}
						break out
					}

				}
			}

		}

		fmt.Fprintf(w, message)
	}
}
