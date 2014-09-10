package app

var path_re = regexp.MustCompile(`^/([0-9a-f]{64})([.][0-9a-z_.-]+)?$`)

type path_t []string

func matchPath(path string) path_t {
  return path_re.FindStringSubmatch(strings.ToLower(path))
}

func (match path_t) hash() string {
  return match[1]
}

func (match path_t) extension() string {
  if len(match) == 2 {
    return ""
  } else {
    return match[2]
  }
}
