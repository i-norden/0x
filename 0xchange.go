package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/i-norden/golimiter"
)

// Exchange structs

// Top level exchange struct
type exchange struct {
	sync.Mutex
	users   []*user
	markets []*market
}

// Struct representing a market for a single pair
type market struct {
	sync.Mutex
	pair      string
	orderbook struct {
		bids []*level
		asks []*level
	}
}

// Struct representing a level on an orderbook
type level struct {
	rate          float64
	totalQuantity float64
	orders        []*limitOrder
}

// Struct representing a limit order
type limitOrder struct {
	time       time.Time
	side       string
	rate       float64
	quantity   float64
	orderid    string
	conditions map[string]interface{}
}

// Struct representing a market order
type marketOrder struct {
	time       time.Time
	side       string
	quantity   float64
	orderid    string
	conditions map[string]interface{}
}

// User structs

// Struct representing a user
type user struct {
	sync.Mutex
	id      string
	account account
}

// Struct representing a user account
type account struct {
	balances     map[string]balances
	orders       map[string]limitOrder // map of orderIDs to the corresponding order
	orderHistory struct {
		limiterOrders []*limitOrder
		marketOrders  []*marketOrder
	}
	balanceHistory struct {
		withdraws []*deposit
		deposits  []*withdraw
	}
}

// Struct representing a user's balances
type balances struct {
	available   float64
	unavailable float64
	total       float64
}

// Management structs

// Struct representing process manager
type management struct {
	sync.Mutex
	activeProcesses map[string]*os.Process
	exchange        *exchange
}

// Struct representing a withdraw
type withdraw struct {
	id       string
	userid   string
	txid     string
	time     time.Time
	asset    string
	quantity float64
	status   string
}

// Struct representing a deposit
type deposit struct {
	id       string
	userid   string
	txid     string
	time     time.Time
	asset    string
	quantity float64
	status   string
}

// Http handler functions for user API endpoints

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

// Function for adding limit order to orderbook
func (m *market) addOrder(order limitOrder) {

}

// Function for filling a market order from an orderbook
func (m *market) fillOrder(order marketOrder) {

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
