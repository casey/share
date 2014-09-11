package app

import "net/http"

type status_t int

func (status status_t) informational() bool { return status >= 100 && status < 200 }
func (status status_t) successful() bool    { return status >= 200 && status < 300 }
func (status status_t) redirection() bool   { return status >= 300 && status < 400 }
func (status status_t) badRequest() bool    { return status >= 400 && status < 500 }
func (status status_t) serverError() bool   { return status >= 500 && status < 600 }

func (status status_t) mustNotIncludeMessageBody(method string) bool {
  return status.informational() ||
    status == http.StatusNoContent ||
    status == http.StatusResetContent ||
    status == http.StatusNotModified ||
    status == http.StatusOK && method == "HEAD"
}

func (status status_t) text() string {
  if text := http.StatusText(status.number()); text != "" {
    return text
  }

  switch status {
  case 420:
    return "Enhance Your Calm"
  case 451:
    return "Unavailable For Legal Reasons"
  case 522:
    return "Unprocessable Entity"
  }

  return "Mystery Status Code"
}

func (status status_t) number() int {
  return int(status)
}
