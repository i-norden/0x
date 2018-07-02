//make it so that if a client closes a websocket connection any ruby processes that are no longer needed for other feeds are exited
//make it so client can choose what period to query at and also which indicator set (out of the 4 premade first)
package main

import (
	"time"
	"log"
	"net/http"
  "sync"
	"os"

  "github.com/gorilla/mux"
  "github.com/i-norden/golimiter"
)


// Exchange structs

type exchange struct {
	users []user
	markets []market
}

type market struct {
	sync.Mutex
	pair string
	orderbook struct {
		bids []level
		asks []level
	}
}

type level struct {
	rate float64
	totalQuantity float64
	orders	[]limitOrder
}

type limitOrder struct {
	time time.Time
	side string
	rate float64
	quantity float64
	orderID string
	conditions map[string]interface{}
}

type marketOrder struct {
	time time.Time
	side string
	quantity float64
	orderID string
	conditions map[string]interface{}
}

// User structs

type user struct {
	sync.Mutex
	account account
}

type account struct {
	id string
	balances map[string]balances
	orders map[string]limitOrder // map of orderIDs to the corresponding order
	orderHistory struct {
		limiterOrders []limitOrder
		marketOrders  []marketOrder
	}
}

type balances struct {
	available float64
	unavailable float64
	total float64
}

// Management structs

type management struct {
	sync.Mutex
	activeProcesses map[string]*os.Process
}


// Http handler functions for API endpoints

func TokenPairs(w http.ResponseWriter, r *http.Request) {

}

func Token(w http.ResponseWriter, r *http.Request) {

}

func Orders(w http.ResponseWriter, r *http.Request) {

}

func OrderHash(w http.ResponseWriter, r *http.Request) {

}

func Orderbook(w http.ResponseWriter, r *http.Request) {

}

func Fees(w http.ResponseWriter, r *http.Request) {

}

func Order(w http.ResponseWriter, r *http.Request) {

}

// Main process

func main() {
  r := mux.NewRouter()
  r.Methods("GET", "POST")
	r.HandleFunc("/v0/token_pairs", TokenPairs)
  r.HandleFunc("/v0/token_pairs/{token}", Token)
	r.HandleFunc("/v0/orders", Orders)
	r.HandleFunc("/v0/order/{orderHash}", OrderHash)
	r.HandleFunc("/v0/orderbook", Orderbook)
  r.HandleFunc("/v0/fees", Fees)
  r.HandleFunc("/v0/order", Order)
	lim := golimiter.Limiter{}
	lim.Rate = 1
	lim.Burst = 5
	err := lim.Init()
	if err != nil {
		log.Fatal("Unable to initiate api limiter")
	}
	for {
		go log.Fatal(http.ListenAndServe(":8080", lim.LimitHTTPHandler(r)))
	}
}
