//make it so that if a client closes a websocket connection any ruby processes that are no longer needed for other feeds are exited
//make it so client can choose what period to query at and also which indicator set (out of the 4 premade first)
package main

import (
	"flag"
	"log"
	"net/http"
	"time"
  "sync"

  "github.com/gorilla/mux"
  "golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

//create
var limiter = rate.NewLimiter(1, 6)

// Create a map to hold the visitor structs for each ip
var visitors = make(map[string]*visitor)
var mtx sync.Mutex

// Create a new rate limiter and add it to the visitors map, using the
// IP address as the key.
func addVisitor(ip string) *rate.Limiter {
    limiter := rate.NewLimiter(2, 5)
    mtx.Lock()
    // Include the current time when creating a new visitor.
    visitors[ip] = &visitor{limiter, time.Now()}
    mtx.Unlock()
    return limiter
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise call the addVisitor function to add a
// new entry to the map.
func getVisitor(ip string) *rate.Limiter {
    mtx.Lock()
    v, exists := visitors[ip]
    if !exists {
        mtx.Unlock()
        return addVisitor(ip)
    }
    // Update the last seen time for the visitor.
    v.lastSeen = time.Now()
    mtx.Unlock()
    return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
    for {
        time.Sleep(time.Minute)
        mtx.Lock()
        for ip, v := range visitors {
            if time.Now().Sub(v.lastSeen) > 3*time.Minute {
                delete(visitors, ip)
            }
        }
        mtx.Unlock()
    }
}

//checks incoming ip against their limit
func limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Call the getVisitor function to retreive the rate limiter for
        // the current user.
        limiter := getVisitor(r.RemoteAddr)
        if limiter.Allow() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

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

var activeProcesses = make(map[string]bool)

func main() {
	flag.Parse()
	log.SetFlags(0)
  r := mux.NewRouter()
  r.Methods("GET", "POST")
	r.HandleFunc("/v0/token_pairs", TokenPairs)
  r.HandleFunc("/v0/token_pairs/{token}", Token)
	r.HandleFunc("/v0/orders", Orders)
	r.HandleFunc("/v0/order/{orderHash}", OrderHash)
	r.HandleFunc("/v0/orderbook", Orderbook)
  r.HandleFunc("/v0/fees", Fees)
  r.HandleFunc("/v0/order", Order)
	go log.Fatal(http.ListenAndServe(":8080", limit(r)))
}

// Run a background goroutine to remove old entries from the visitors map.
func init() {
    go cleanupVisitors()
}
