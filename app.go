package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sony/sonyflake"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})
var roomidSet = make(map[uint64]bool)
var nicknameSet = make(map[uint64][]string)

func main() {
	// Http Handlers
	http.Handle("/", http.FileServer(http.Dir("views")))

	http.HandleFunc("/create", CreateRoom)
	http.HandleFunc("/join", JoinRoom)

	fmt.Println("Listening to port 9797...")
	http.ListenAndServe(":9797", nil)

	// Socekt Handlers
	// server, err := socketio.NewServer(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// server.OnConnect("/", SocketConnect)
}

/*##################### Handlers #####################*/

/*CreateRoom generates new GUID, creates a socket room using it*/
func CreateRoom(res http.ResponseWriter, req *http.Request) {
	type reqStruct struct {
		Nickname string `json:"nickname"`
	}
	var data reqStruct
	GetBodyData(req, &data)
	nickname := data.Nickname

	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}

	roomidSet[id] = true

	names, ok := nicknameSet[id]
	if ok {
		nicknameSet[id] = append(names, nickname)
	} else {
		nicknameSet[id] = []string{nickname}
	}

	fmt.Fprintf(res, strconv.FormatUint(id, 10))
}

/*JoinRoom checks if the recieved id exists, then connect the socket to the room*/
func JoinRoom(res http.ResponseWriter, req *http.Request) {
	type reqStruct struct {
		Nickname string `json:"nickname"`
		Roomcode string `json:"roomCode"`
	}
	var data reqStruct
	GetBodyData(req, &data)
	nickname := data.Nickname
	roomcode := data.Roomcode
	code, _ := strconv.ParseUint(roomcode, 10, 64)

	_, ok := roomidSet[code]
	// Check if room exists
	if ok {
		users := nicknameSet[code]
		if NameExists(users, nickname) {
			fmt.Fprintf(res, "Name exists")
		} else {
			nicknameSet[code] = append(users, nickname)
			fmt.Fprintf(res, "Available")
		}
	} else {
		fmt.Fprintf(res, "Not available")
	}
}

/*##################### Helper #####################*/

/*CheckError used by GetBodyData to handle errors*/
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

/*GetBodyData parses http body for json data*/
func GetBodyData(req *http.Request, class interface{}) {
	err := json.NewDecoder(req.Body).Decode(&class)
	CheckError(err)
}

/*NameExists checks if an item exists in an array*/
func NameExists(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
