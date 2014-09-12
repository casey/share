package app

import "errors"

type response_t struct {
  err         error
  status      status_t
  body        string
  contentType string
}

func (response response_t) finish() {
  panic(response)
}

func empty(status status_t) {
  response_t{nil, status, "", ""}.finish()
}

func full(status status_t, body string, contentType string) {
  response_t{nil, status, body, contentType}.finish()
}

func ensure(condition bool, status status_t) {
  if !condition {
    response_t{errors.New("ensure condition false"), status, "", ""}.finish()
  }
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
