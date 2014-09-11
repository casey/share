package app

import "appengine"
import "fmt"
import "net/http"

const maximumContentLength = 128

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  body := (*string)(nil)
  status := status_t(http.StatusInternalServerError)
  contentType := `text/plain; charset="utf-8"`

  defer func() {
    relic := recover()
    
    if relic == nil {
      c.Errorf("handler: completed without response")
    } else {
      c.Infof("handler: completed normally")
    }

    w.Header().Set("License", "Anyone may do anything with this.")
    w.Header().Set("Warranty", `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`)
    w.Header().Set("Content-Type", contentType)

    w.WriteHeader(status.number())

    if status.mustNotIncludeMessageBody(r.Method) {
      fmt.Fprint(w, "\n")
    } else if body == nil {
      fmt.Fprintf(w, "%v %v\n", status.number(), status.text())
    } else {
      fmt.Fprintf(w, "%v", *body)
    }
  }()

  ensure := func(condition bool, errorCode int) {
    if !condition {
      status = status_t(errorCode)
      panic("ensure condition false")
    }
  }

  check := func(e error) {
    if e != nil {
      panic(e)
    }
  }

  empty := func(s status_t) {
    status = s
    panic("response without body")
  }

  corporeal := func(s status_t, b string, ct string) {
    status = s
    body = new(string)
    *body = b
    contentType = ct
    panic("response with body")
  }

  get := r.Method == "GET"
  put := r.Method == "PUT"
  ensure(put || get, http.StatusMethodNotAllowed)
  match := matchPath(r.URL.Path)
  ensure(match != nil, http.StatusForbidden)

  if get {
    pointer, e := fetch(c, match.hash())
    check(e)
    ensure(pointer != nil, http.StatusNotFound)
    corporeal(http.StatusOK, string(*pointer), "application/octet-stream")
  } else {
    ensure(match.extension() == "", http.StatusForbidden)
    ensure(r.ContentLength >= 0, http.StatusLengthRequired)
    ensure(r.ContentLength <= maximumContentLength, http.StatusRequestEntityTooLarge)
    buffer := make([]byte, r.ContentLength)
    n, e := r.Body.Read(buffer)
    check(e)
    ensure(int64(n) == r.ContentLength, http.StatusInternalServerError)
    ensure(hashOK(match.hash(), buffer), http.StatusForbidden)
    p, e := published(c, match.hash());
    check(e)
    if p {
      empty(http.StatusOK)
    } else {
      check(publish(c, match.hash(), buffer))
      empty(http.StatusCreated)
    }
  }
}
