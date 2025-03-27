# nls
Yet another easy to use i18n localization generator.

The tool `nls` loads all messages (in YAML files) in all available languages and generates a Go package with Go constants per message and a Localizer with lookup and replace functions.

## install

    go install github.com/emicklei/cmd/nls@latest

## tool usage
Initialy, you start with a folder (e.g. `messages/en`) in your project with an empty `messages.yaml` file.
After adding a message key, you run `go generate` to (re)generate the Go package and update all other languages with missing keys.

## package usage
```go
package main

import (
	"fmt"

	"github.com/emicklei/nls/example/nls"
	"golang.org/x/text/language"
)

//go:generate nls -dir messages -pkg nls
func main() {
	loc := nls.New(language.Dutch.String(), language.English.String())

	fmt.Println(loc.Get(nls.M_hello))
	fmt.Println(loc.Get(nls.M_world))
	fmt.Println(loc.Get(nls.M_sky, "Sky"))
	fmt.Println(loc.Replaced(nls.M_sea, map[string]string{"name": "Noord"}))
	fmt.Println(loc.Replaced(nls.M_cats, map[string]any{"count": 3}))
	fmt.Println(loc.Format(nls.M_cats, "count", 1))
}
```
Outputs
```
Hallo
Wereld
Sky
blauwe zee
3 katten
1 kat
```

## acknowledgements

Making of this package is inspired by the (inactive) [go-localize](https://github.com/m1/go-localize) package.