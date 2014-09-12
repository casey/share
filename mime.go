package app

import "mime"
import "net/http"

func getContentType(extension string, data []byte) string {
  guess := ""

  if extension == ".sniff" {
    guess = http.DetectContentType(data)
  } else if extension != "" {
    guess = mime.TypeByExtension(extension)
  }

  if guess == "" {
    return "application/octet-stream"
  }

  return guess
}
