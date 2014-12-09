package main

import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "encoding/json"
import "net/http"
import "strings"
import "net/url"
import "fmt"

type Response struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

var Version = "0.0.1"

const Usage = `
  Usage:
    go-search <query>...
    go-search -h | --help
    go-search --version

  Options:
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	query := strings.Join(args["<query>"].([]string), " ")

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

	println()
	for _, pkg := range body.Results {
		fmt.Printf("  \033[1m%s\033[m\n", strip(pkg.Path))
		if pkg.Synopsis == "" {
			fmt.Printf("  %s\n", "no description")
		} else {
			fmt.Printf("  %s\n", pkg.Synopsis)
		}
		fmt.Printf("  godoc.org/pkg/%s\n\n", pkg.Path)
	}
	println()
}

func strip(path string) string {
	return strings.Replace(path, "github.com/", "", 1)
}
