package app

import "appengine"
import "fmt"
import "net/http"
import "regexp"

var path_re = regexp.MustCompile(`^/([0-9a-f]{64})([.][0-9a-z_.-]+)?$`)

const maximumContentLength = 128

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  defer func() {
    c := appengine.NewContext(r)
    relic := recover()
    response, ok := relic.(response_t)

    if !ok {
      c.Errorf("handler: relic was not response: %v", relic)
      response.status = http.StatusInternalServerError
    }

    if response.status.mustNotIncludeMessageBody(r.Method) {
      response.body = "\n"
    } else if response.body == "" && response.contentType == "" {
      response.body = fmt.Sprintf("%v %v\n", response.status.number(), response.status.text())
    }

    if response.contentType == "" {
      response.contentType = `text/plain; charset="utf-8"`
    }

    w.Header().Set("License", "Anyone may do anything with this.")
    w.Header().Set("Warranty", `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`)
    w.Header().Set("Content-Type", response.contentType)
    w.WriteHeader(response.status.number())
    w.Write([]byte(response.body))
  }()

  switch r.Method {
  case "GET": get(r)
  case "PUT": put(r)
  default:     empty(http.StatusMethodNotAllowed)
  }
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  match := matchPath(r.URL.Path)
  ensure(match != nil, http.StatusForbidden)
  pointer, e := fetch(c, match.hash())
  check(e)
  ensure(pointer != nil, http.StatusNotFound)
  full(http.StatusOK, string(*pointer), "application/octet-stream")
}

func put(r *http.Request) {
  c := appengine.NewContext(r)
  match := matchPath(r.URL.Path)
  ensure(match != nil, http.StatusForbidden)
  ensure(match.extension() == "", http.StatusForbidden)
  ensure(r.ContentLength >= 0, http.StatusLengthRequired)
  ensure(r.ContentLength <= maximumContentLength, http.StatusRequestEntityTooLarge)
  buffer := make([]byte, r.ContentLength)
  n, e := r.Body.Read(buffer)
  check(e)
  ensure(int64(n) == r.ContentLength, http.StatusInternalServerError)
  ensure(hashOK(match.hash(), buffer), http.StatusForbidden)
  s, e := shared(c, match.hash())
  check(e)
  ensure(s, http.StatusOK)
  check(share(c, match.hash(), buffer))
  empty(http.StatusCreated)
}
