package app

import "strings"

type path_t []string

func matchPath(path string) path_t {
  return path_re.FindStringSubmatch(strings.ToLower(path))
}

func (match path_t) hash() string {
  return match[1]
}

func (match path_t) extension() string {
  if len(match) == 3 {
    return match[2]
  }
  return ""
}
