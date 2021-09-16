package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

var board [10][10]int
var ships [3]ship
var takenSquares [10][10]bool

func main() {
	size := 10

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

	//fmt.Print(ships)

	/* port := os.Args[1]
	=	fmt.Println("Port = ", port)

		listenPort := ":" + port */

	//mux := http.NewServeMux()

	http.HandleFunc("/board", BoardHandler)
	http.HandleFunc("/boats", BoatsHandler)
	http.HandleFunc("/hit", HitHandler)

	//WrappedMux := RunSomeCode(mux)
	go http.ListenAndServe(":9001", nil)

	scanner := bufio.NewScanner(os.Stdin)

	var addresses []string
out:
	for true {
		scanner.Scan()
		switch scanner.Text() {
		case "connect":
			fmt.Println("Enter an address to connect")
			scanner.Scan()
			addresses = append(addresses, scanner.Text())

		case "attack":
			fmt.Println("Choose a player to attack ;")
			for i := 1; i <= len(addresses); i++ {
				fmt.Println(i, "=>", addresses[i-1])
			}
			scanner.Scan()
			num, _ := strconv.Atoi(scanner.Text())
			apiUrl := "http://" + addresses[num-1]

			response, _ := http.Get(apiUrl + "/board")
			board, _ := ioutil.ReadAll(response.Body)
			sb := string(board)
			fmt.Println(sb)

			resource := "/hit"
			data := url.Values{}
			fmt.Print("X : ")
			scanner.Scan()
			x := scanner.Text()
			fmt.Print("Y : ")
			y := scanner.Text()
			scanner.Scan()
			data.Set("x", x)
			data.Set("y", y)
			u, _ := url.ParseRequestURI(apiUrl)
			u.Path = resource
			urlStr := u.String() // "https://api.com/user/"
			client := &http.Client{}
			r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
			r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			response2, _ := client.Do(r)
			resp, _ := ioutil.ReadAll(response2.Body)
			sb2 := string(resp)
			fmt.Println(sb2)

		case "test":
			response, _ := http.Get("http://localhost:9001/board")
			data, _ := ioutil.ReadAll(response.Body)
			sb := string(data)
			fmt.Println(sb)

		case "exit":
			break out
		}
	}

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
