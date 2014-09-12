package app

import "mime"
import "net/http"

func getContentType(extension string, data []byte) string {
  guess := ""

  if extension == ".sniff" {
    print("sniffing")
    guess = http.DetectContentType(data)
  } else if extension != "" {
    print("guessing by extension")
    guess = mime.TypeByExtension(extension)
  }

  if guess == "" {
    return "application/octet-stream"
  }

  return guess
}
