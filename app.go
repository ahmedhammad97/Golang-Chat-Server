package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sony/sonyflake"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})
var roomidSet = make(map[uint64]bool)
var nicknameSet = make(map[uint64][]string)

func main() {
	// Socket Handlers
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", SocketConnect)
	server.OnEvent("/", "initialConnection", InitConn)
	server.OnEvent("/", "chatMessage", ChatMessage)
	server.OnEvent("/", "typing", Typing)

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	// Http Handlers
	http.Handle("/", http.FileServer(http.Dir("views")))
	http.Handle("/socket", server)

	http.HandleFunc("/create", CreateRoom)
	http.HandleFunc("/join", JoinRoom)

	fmt.Println("Listening to port 9797...")
	http.ListenAndServe(":9797", nil)
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

/*SocketConnect handles first socket connection*/
func SocketConnect(so socketio.Conn) error {
	so.Emit("initialConnection")
	return nil
}

/*InitConn allow socket to join a room*/
func InitConn(so socketio.Conn, msg string) {
	so.Join(msg)
}

/*ChatMessage broadcasts message to other users*/
func ChatMessage(so socketio.Conn, msg string) {
	so.Emit("chatMessage", msg)
}

/*Typing notify other users that this user is typing*/
func Typing(so socketio.Conn, msg string) {
	so.Emit("typing", msg)
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
