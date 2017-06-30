package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-wordwrap"
	gopen "github.com/petermbenjamin/go-open"
	"github.com/tj/docopt"
)

// Response represents the response to be returned
type Response struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

// Version is the package version
var Version = "0.0.1"

// Usage is the package usage information
const Usage = `
  Usage:
    go-search <query>... [--top] [--count n] [--open]
    go-search -h | --help
    go-search --version

  Options:
    -n, --count n    number of results [default: -1]
    -t, --top        top-level packages only
    -o, --open       open godoc.org search results in default browser
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	n, err := strconv.ParseInt(args["--count"].(string), 10, 32)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	query := strings.Join(args["<query>"].([]string), " ")
	top := args["--top"].(bool)

	res, err := http.Get("https://api.godoc.org/search?q=" + url.QueryEscape(query))
	if err != nil {
		log.Fatalf("request failed: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("request error: %s", http.StatusText(res.StatusCode))
	}

	var body Response
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatalf("error parsing response: %s", err)
	}

	if open := args["--open"].(bool); open {
		gopen.Open("https://godoc.org/?q=" + url.QueryEscape(query))
		os.Exit(0)
	}

	if n > 0 {
		body.Results = body.Results[:n]
	}

	println()
	for _, pkg := range body.Results {
		if top && subpackage(pkg.Path) {
			continue
		}

		fmt.Printf("  \033[1m%s\033[m\n", pkg.Path)
		fmt.Printf("  godoc.org/pkg/%s\n", pkg.Path)
		fmt.Printf("  %s\n", description(pkg.Synopsis))
		fmt.Printf("\n")
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

	s = wordwrap.WrapString(s, 60)
	s = strings.Replace(s, "\n", "\n  ", -1)
	return s
}
