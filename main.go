package main

import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "encoding/json"
import "net/http"
import "strings"
import "net/url"
import "fmt"
import "github.com/mattn/go-colorable"

type Response struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

var Version = "0.0.1"

const Usage = `
  Usage:
    go-search <query>... [--top]
    go-search -h | --help
    go-search --version

  Options:
    -t, --top        top-level packages only
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	query := strings.Join(args["<query>"].([]string), " ")
	top := args["--top"].(bool)

	res, err := http.Get("http://api.godoc.org/search?q=" + url.QueryEscape(query))
	if err != nil {
		log.Critical("request failed: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Critical("request error: %s", http.StatusText(res.StatusCode))
	}

	var body Response
	log.Check(json.NewDecoder(res.Body).Decode(&body))

	out := colorable.NewColorableStdout()

	fmt.Println()
	for _, pkg := range body.Results {
		if top && subpackage(pkg.Path) {
			continue
		}

		fmt.Fprintf(out, "  \033[1m%s\033[m\n", strip(pkg.Path))
		fmt.Fprintf(out, "  %s\n", description(pkg.Synopsis))
		fmt.Fprintf(out, "  godoc.org/pkg/%s\n\n", pkg.Path)
	}
	fmt.Println()
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
