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
	fmt.Println(loc.Format(nls.M_sea, "name", "Noord"))
	fmt.Println(loc.Format(nls.M_cats, "count", 3))
	fmt.Println(loc.Format(nls.M_cats, "count", 1))
}
