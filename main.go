package app

import "appengine"
import "net/http"
import "regexp"

import . "flotilla"

var   path_re              = regexp.MustCompile(`^/(?P<hash>[0-9a-f]{64})(?P<extension>[.][0-9a-z_.-]+)?$`)
const maximumContentLength = 128
const license              = "Anyone may do anything with this."
const warranty             = `"AS IS" WITH NO WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  defer Respond(w, r)
  switch r.Method {
  case "GET": get(r)
  case "PUT": put(r)
  default:    Status(http.StatusMethodNotAllowed)
  }
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  match := Components(path_re, r.URL.Path)
  Ensure(match != nil, http.StatusForbidden)
  pointer, e := fetch(c, match["hash"])
  Check(e)
  Ensure(pointer != nil, http.StatusNotFound)
  Body(http.StatusOK, string(*pointer), sniff(match["extension"], *pointer))
}

func put(r *http.Request) {
  c := appengine.NewContext(r)
  match := Components(path_re, r.URL.Path)
  Ensure(r.Header.Get("License") == license, http.StatusForbidden)
  Ensure(match != nil, http.StatusForbidden)
  Ensure(match["extension"] == "", http.StatusForbidden)
  Ensure(r.ContentLength >= 0, http.StatusLengthRequired)
  Ensure(r.ContentLength <= maximumContentLength, http.StatusRequestEntityTooLarge)
  buffer := make([]byte, r.ContentLength)
  n, e := r.Body.Read(buffer)
  Check(e)
  Ensure(int64(n) == r.ContentLength, http.StatusInternalServerError)
  Ensure(hashOK(match["hash"], buffer), http.StatusForbidden)
  shared, e := shared(c, match["hash"])
  Check(e)
  if shared {
    Status(http.StatusOK)
  }
  Check(share(c, match["hash"], buffer))
  Status(http.StatusCreated)
}
