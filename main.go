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

    if response.why != "" {
      c.Infof("handler: why: %v", response.why)
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
    w.Header().Set("Warranty", `"AS IS" WITH NO WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", response.contentType)
    w.WriteHeader(response.status.number())
    w.Write([]byte(response.body))
  }()

  switch r.Method {
  case "GET": get(r)
  case "PUT": put(r)
  default:    status(http.StatusMethodNotAllowed)
  }
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  match := matchPath(r.URL.Path)
  ensure(match != nil, http.StatusForbidden)
  pointer, e := fetch(c, match.hash())
  check(e)
  ensure(pointer != nil, http.StatusNotFound)
  body(http.StatusOK, string(*pointer), "application/octet-stream")
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
  shared, e := shared(c, match.hash())
  check(e)
  if shared {
    status(http.StatusOK)
  }
  check(share(c, match.hash(), buffer))
  status(http.StatusCreated)
}
