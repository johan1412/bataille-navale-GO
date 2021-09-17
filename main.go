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

type player struct {
	address string
	isAlive bool
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
	http.HandleFunc("/message", MessageHandler)

	port := os.Args[1]
	//WrappedMux := RunSomeCode(mux)
	go http.ListenAndServe(":"+port, nil)

	scanner := bufio.NewScanner(os.Stdin)

	var players []player

	fmt.Print("Username : ")
	scanner.Scan()
	username := scanner.Text()

out:
	for true {

		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(100)
		if random == 42 {
			tsunami()
		}

		fmt.Print("\nAvailable commands:\n\n- connect : Connect to a player\n- attack : Attack one of the players you're connected ton\n- attack-special : Attack special one of the players you're connected to\n- message : Message one of the players you're connected to\n- exit : Exit the game\n\n")
		scanner.Scan()
		switch scanner.Text() {
		case "message":
			if len(players) == 0 {
				fmt.Print("\nYou are not connected to any player\n\n")
			} else {
				fmt.Print("\nChoose a player to message :\n\n")
				for i := 1; i <= len(players); i++ {
					fmt.Println(i, "=>", players[i-1].address)
				}
				scanner.Scan()
				num, _ := strconv.Atoi(scanner.Text())
				fmt.Print("\nEnter message :\n\n")
				scanner.Scan()
				message := scanner.Text()
				resource := "/message"
				data := url.Values{}
				data.Set("username", username)
				data.Set("message", message)
				u, _ := url.ParseRequestURI(players[num-1].address)
				u.Path = resource
				urlStr := u.String()
				client := &http.Client{}
				r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
				r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
				client.Do(r)

			}

		case "connect":
			fmt.Print("\nEnter an address to connect ( example : http://127.0.0.1:8000 ) \n\n")
			scanner.Scan()
			var newPlayer player
			newPlayer.address = scanner.Text()
			newPlayer.isAlive = true
			players = append(players, newPlayer)

		case "attack":
			if getRemainingShips(ships) == 0 {
				fmt.Print("\nYou cannot attack, you have 0 remaining ships\n\n")
			} else {
				if len(players) == 0 {
					fmt.Print("\nYou are not connected to any player\n\n")
				} else {
					fmt.Print("\nChoose a player to attack :\n\n")
					for i := 1; i <= len(players); i++ {
						if players[i-1].isAlive {
							fmt.Println(i, "=>", players[i-1].address)
						}
					}
					scanner.Scan()
					num, _ := strconv.Atoi(scanner.Text())
					apiUrl := players[num-1].address

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
					scanner.Scan()
					y := scanner.Text()
					data.Set("x", x)
					data.Set("y", y)
					u, _ := url.ParseRequestURI(apiUrl)
					u.Path = resource
					urlStr := u.String()
					client := &http.Client{}
					r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
					r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
					response2, _ := client.Do(r)
					resp, _ := ioutil.ReadAll(response2.Body)
					sb2 := string(resp)
					fmt.Print("\n", sb2, "\n\n")
					response3, _ := http.Get(players[num-1].address + "/boats")
					data3, _ := ioutil.ReadAll(response3.Body)
					sb3 := string(data3)
					fmt.Println("Remaining ships :", sb3)

					hasWon := true

					for i := 0; i < len(players); i++ {

						response, _ := http.Get(players[i].address + "/boats")
						data, _ := ioutil.ReadAll(response.Body)
						sb := string(data)
						nb := 0
						nb, _ = strconv.Atoi(sb)
						if nb == 0 {
							players[i].isAlive = false
						}
						if players[i].isAlive {
							hasWon = false
						}
					}
					if hasWon {
						fmt.Print("\nYou won !\n")
						break out
					}

				}
			}
		case "attack-special":
			if getRemainingShips(ships) == 0 {
				fmt.Print("\nYou cannot attack, you have 0 remaining ships\n\n")
			} else {
				if len(players) == 0 {
					fmt.Print("\nYou are not connected to any player\n\n")
				} else {
					fmt.Print("\nChoose a player to attack :\n\n")
					for i := 1; i <= len(players); i++ {
						if players[i-1].isAlive {
							fmt.Println(i, "=>", players[i-1].address)
						}
					}
					scanner.Scan()
					num, _ := strconv.Atoi(scanner.Text())
					apiUrl := players[num-1].address

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
					scanner.Scan()
					y := scanner.Text()
					data.Set("x", x)
					data.Set("y", y)

					u, _ := url.ParseRequestURI(apiUrl)
					u.Path = resource
					urlStr := u.String()
					client := &http.Client{}
					r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
					r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
					response2, _ := client.Do(r)

					resp, _ := ioutil.ReadAll(response2.Body)
					sb2 := string(resp)
					fmt.Print("\n", sb2, "\n\n")

					x_int, _ := strconv.Atoi(x)
					y_int, _ := strconv.Atoi(y)

					x2 := x_int + 1 // down
					x3 := x_int - 1 // up

					y2 := y_int + 1 // rigth
					y3 := y_int - 1 // left

					// attack downn
					if x2 > 0 && x2 < 10 {

						data2 := url.Values{}
						data2.Set("x", strconv.Itoa(x2))
						data2.Set("y", strconv.Itoa(y_int))

						u1, _ := url.ParseRequestURI(apiUrl)
						u1.Path = resource
						urlStr1 := u.String()
						client1 := &http.Client{}
						r1, _ := http.NewRequest(http.MethodPost, urlStr1, strings.NewReader(data2.Encode())) // URL-encoded payload
						r1.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
						r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						r1.Header.Add("Content-Length", strconv.Itoa(len(data2.Encode())))
						response22, _ := client1.Do(r1)

						resp2, _ := ioutil.ReadAll(response22.Body)
						sb22 := string(resp2)
						fmt.Print("\n", sb22, "\n\n")
					}

					// attack up
					if x3 > 0 && x3 < 10 {
						data2 := url.Values{}
						data2.Set("x", strconv.Itoa(x3))
						data2.Set("y", strconv.Itoa(y_int))

						u1, _ := url.ParseRequestURI(apiUrl)
						u1.Path = resource
						urlStr1 := u.String()
						client1 := &http.Client{}
						r1, _ := http.NewRequest(http.MethodPost, urlStr1, strings.NewReader(data2.Encode())) // URL-encoded payload
						r1.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
						r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						r1.Header.Add("Content-Length", strconv.Itoa(len(data2.Encode())))
						response22, _ := client1.Do(r1)

						resp2, _ := ioutil.ReadAll(response22.Body)
						sb22 := string(resp2)
						fmt.Print("\n", sb22, "\n\n")
					}

					// attack right
					if y2 > 0 && y2 < 10 {

						data2 := url.Values{}
						data2.Set("x", strconv.Itoa(x_int))
						data2.Set("y", strconv.Itoa(y2))

						u1, _ := url.ParseRequestURI(apiUrl)
						u1.Path = resource
						urlStr1 := u.String()
						client1 := &http.Client{}
						r1, _ := http.NewRequest(http.MethodPost, urlStr1, strings.NewReader(data2.Encode())) // URL-encoded payload
						r1.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
						r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						r1.Header.Add("Content-Length", strconv.Itoa(len(data2.Encode())))
						response22, _ := client1.Do(r1)

						resp2, _ := ioutil.ReadAll(response22.Body)
						sb22 := string(resp2)
						fmt.Print("\n", sb22, "\n\n")
					}

					// attack left
					if y3 > 0 && y3 < 10 {

						data2 := url.Values{}
						data2.Set("x", strconv.Itoa(x_int))
						data2.Set("y", strconv.Itoa(y3))

						u1, _ := url.ParseRequestURI(apiUrl)
						u1.Path = resource
						urlStr1 := u.String()
						client1 := &http.Client{}
						r1, _ := http.NewRequest(http.MethodPost, urlStr1, strings.NewReader(data2.Encode())) // URL-encoded payload
						r1.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
						r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						r1.Header.Add("Content-Length", strconv.Itoa(len(data2.Encode())))
						response22, _ := client1.Do(r1)

						resp2, _ := ioutil.ReadAll(response22.Body)
						sb22 := string(resp2)
						fmt.Print("\n", sb22, "\n\n")
					}

					response3, _ := http.Get(players[num-1].address + "/boats")
					data3, _ := ioutil.ReadAll(response3.Body)
					sb3 := string(data3)
					fmt.Println("Remaining ships :", sb3)

					hasWon := true

					for i := 0; i < len(players); i++ {

						response, _ := http.Get(players[i].address + "/boats")
						data, _ := ioutil.ReadAll(response.Body)
						sb := string(data)
						nb := 0
						nb, _ = strconv.Atoi(sb)
						if nb == 0 {
							players[i].isAlive = false
						}
						if players[i].isAlive {
							hasWon = false
						}
					}
					if hasWon {
						fmt.Print("\nYou won !\n")
						break out
					}

				}
			}

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

func getRemainingShips(ships [3]ship) int {
	remainingShips := 0
	for i := 0; i < len(ships); i++ {
		if !ships[i].isSunk() {
			remainingShips++
		}
	}
	return remainingShips
}

func tsunami() {
	fmt.Println("You have been struck by a tsunami !!! You lost.")
	for i := 0; i < len(ships); i++ {
		for j := 0; j < len(ships[i].squares); j++ {
			ships[i].squares[j].hit = true
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			board[i][j] = 3
		}
	}
}
