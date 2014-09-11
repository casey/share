package app

import "appengine"
import "fmt"
import "net/http"
import "crypto/sha256"
import "encoding/hex"

const maximumContentLength = 128

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  status := status_t(http.StatusInternalServerError)
  var body *string
  headers := make(map[string]string)
  headers["License"] = "Anyone may do anything with this."
  headers["Warranty"] = `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`
  headers["Content-Type"] = `text/plain; charset="utf-8"`

  defer func() {
    if e := recover(); e != nil {
      c.Errorf("handler: recovered from panic: %v", e)
    }

    for name, value := range headers {
      w.Header().Set(name, value)
    }

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
      status = http.StatusInternalServerError
      panic(e)
    }
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
    body = new(string)
    *body = string(*pointer)
  } else {
    ensure(match.extension() == "", http.StatusForbidden)
    ensure(r.ContentLength >= 0, http.StatusLengthRequired)
    ensure(r.ContentLength <= maximumContentLength, http.StatusRequestEntityTooLarge)
    buffer := make([]byte, r.ContentLength)
    n, e := r.Body.Read(buffer)
    check(e)
    ensure(int64(n) == r.ContentLength, http.StatusInternalServerError)

    sha := sha256.New()
    sha.Write(buffer)
    sum := sha.Sum(nil)
    calculatedHash := hex.EncodeToString(sum)
    ensure(calculatedHash == match.hash(), http.StatusForbidden)

    p, e := published(c, match.hash())
    check(e)
    if p {
      status = http.StatusOK
    } else {
      status = http.StatusCreated
    }
  }
}
