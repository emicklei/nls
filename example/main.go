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
	fmt.Println(loc.Replaced(nls.M_sea, map[string]any{"name": "Noord"}))
	fmt.Println(loc.Replaced(nls.M_cats, map[string]any{"count": 3}))
	fmt.Println(loc.Replaced(nls.M_cats, map[string]any{"count": 1}))
}
