//https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
package main

import (
  "fmt"
  "strings"
  "path"
  "net/http"
  "strconv"
)

type App struct {
    // We could use http.Handler as a type here; using the specific type has
    // the advantage that static analysis tools can link directly from
    // h.UserHandler.ServeHTTP to the correct definition. The disadvantage is
    // that we have slightly stronger coupling. Do the tradeoff yourself.
    UserHandler *UserHandler
}

type ProfileHandler struct {
}

type UserHandler struct{
    ProfileHandler *ProfileHandler
}

func (u *UserHandler) HandleGet(id int) {

}

func (u *UserHandler) HandlePost(id int) {

}

func ShiftPath(p string) (head, tail string) {
    p = path.Clean("/" + p)
    i := strings.Index(p[1:], "/") + 1
    if i <= 0 {
        return p[1:], "/"
    }
    return p[1:i], p[i:]
}

func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    var head string
    head, req.URL.Path = ShiftPath(req.URL.Path)
    if head == "user" {
        h.UserHandler.ServeHTTP(res, req)
        return
    }
    http.Error(res, "Not Found", http.StatusNotFound)
}

func (h *UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    var head string
    head, req.URL.Path = ShiftPath(req.URL.Path)
    id, err := strconv.Atoi(head)
    if err != nil {
        http.Error(res, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
        return
    }

    if req.URL.Path != "/" {
        head, tail := ShiftPath(req.URL.Path)
        fmt.Println(tail)
        switch head {
        case "profile":
            // We can't just make ProfileHandler an http.Handler; it needs the
            // user id. Let's instead…
            h.ProfileHandler.Handler(id).ServeHTTP(res, req)
        case "account":
            // Left as an exercise to the reader.
        default:
            http.Error(res, "Not Found", http.StatusNotFound)
        }
        return
    }

    switch req.Method {
    case "GET":
        h.HandleGet(id)
    case "POST":
        h.HandlePost(id)
    default:
        http.Error(res, "Only GET and POST are allowed", http.StatusMethodNotAllowed)
    }
}

func (h *ProfileHandler) Handler(id int) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        // Do whatever
    })
}

func main() {
    a := &App{
        UserHandler: new(UserHandler),
    }
    http.ListenAndServe(":8000", a)
}
