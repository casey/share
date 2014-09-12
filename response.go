package app

import "fmt"
import "runtime/debug"
import "net/http"
import "appengine"

type response_t struct {
  why         string
  status      status_t
  body        string
  contentType string
}

func (response response_t) finish() {
  if appengine.IsDevAppServer() {
    debug.PrintStack()
  }
  panic(response)
}

func status(status status_t) {
  response_t{"status only response", status, "", ""}.finish()
}

func body(status status_t, body string, contentType string) {
  response_t{"response with body", status, body, contentType}.finish()
}

func ensure(condition bool, status status_t) {
  if !condition {
    response_t{"ensure condition false", status, "", ""}.finish()
  }
}

func check(e error) {
  if e != nil {
    response_t{fmt.Sprintf("check failed: %v", e), http.StatusInternalServerError, "", ""}.finish()
  }
}
