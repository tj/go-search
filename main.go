package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/segmentio/go-log"
	"github.com/tj/docopt"
)

type Response struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

var Version = "0.0.1"

const Usage = `
  Usage:
    go-search <query>... [--top] [--count n]
    go-search -h | --help
    go-search --version

  Options:
    -n, --count n    number of results [default: -1]
    -t, --top        top-level packages only
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	n, err := strconv.ParseInt(args["--count"].(string), 10, 32)
	log.Check(err)

	query := strings.Join(args["<query>"].([]string), " ")
	top := args["--top"].(bool)

	res, err := http.Get("http://api.godoc.org/search?q=" + url.QueryEscape(query))
	if err != nil {
		log.Fatalf("request failed: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("request error: %s", http.StatusText(res.StatusCode))
	}

	var body Response
	log.Check(json.NewDecoder(res.Body).Decode(&body))

	if n > 0 {
		body.Results = body.Results[:n]
	}

	println()
	for _, pkg := range body.Results {
		if top && subpackage(pkg.Path) {
			continue
		}

		fmt.Printf("  \033[1m%s\033[m\n", strip(pkg.Path))
		fmt.Printf("  \033[36mdsc:\033[0m %s\n", description(pkg.Synopsis))
		fmt.Printf("  \033[36mdoc:\033[0m godoc.org/pkg/%s\n", pkg.Path)
		fmt.Printf("  \033[36mpkg:\033[0m %s\n\n", pkg.Path)
	}
	println()
}

func subpackage(s string) bool {
	switch {
	case strings.HasPrefix(s, "github.com"):
		return strings.Count(s, "/") > 3
	default:
		return false
	}
}

func description(s string) string {
	if s == "" {
		return "no description"
	}

	return s
}

func strip(s string) string {
	return strings.Replace(s, "github.com/", "", 1)
}
