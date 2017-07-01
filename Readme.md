# go-search

 [Godoc.org](http://godoc.org) via the command-line.

 ![godoc command-line search](https://dl.dropboxusercontent.com/u/6396913/go-search/Screen%20Shot%202014-12-08%20at%207.25.59%20PM.png)

## Installation

Via `go-get` or the binary [releases](https://github.com/tj/go-search/releases).

```
$ go get github.com/tj/go-search
```

## Usage

- Add an alias: `alias gos=go-search`
- Help text:
    ```sh
    $ go-search --help
    Usage:
        go-search <query>... [--top] [--count n] [--open]
        go-search -h | --help
        go-search --version

      Options:
        -n, --count n    number of results [default: 5]
        -t, --top        top-level packages only
        -o, --open       open godoc.org search results in default browser
        -h, --help       output help information
        -v, --version    output version
    ```
- Examples:
    ```sh
    $ go-search UUID

        github.com/pborman/uuid
        godoc.org/pkg/github.com/pborman/uuid
        The uuid package generates and inspects UUIDs.

        github.com/satori/go.uuid
        godoc.org/pkg/github.com/satori/go.uuid
        Package uuid provides implementation of Universally Unique
        Identifier (UUID).

        github.com/nu7hatch/gouuid
        godoc.org/pkg/github.com/nu7hatch/gouuid
        This package provides immutable UUID structs and the
        functions NewV3, NewV4, NewV5 and Parse() for generating
        versions 3, 4 and 5 UUIDs as specified in RFC 4122.

        github.com/twinj/uuid
        godoc.org/pkg/github.com/twinj/uuid
        This package provides RFC4122 and DCE 1.1 UUIDs.

        github.com/docker/distribution/uuid
        godoc.org/pkg/github.com/docker/distribution/uuid
        Package uuid provides simple UUID generation.
    ```
    ```sh
    $ go-search UUID --open
    # opens godoc.org search results in default browser
    ```

## Changelog

### 0.0.3

- Default search results limit in terminal to 5. Users may still control the limit with `-n`.
- Minor fixes.

### 0.0.2

- Feature: open search results on godoc.org with `--open`.

### 0.0.1

- Initial release.

## License

MIT