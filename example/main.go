package main

import (
	"fmt"

	//lint:ignore ST1001 less verbose
	. "github.com/emicklei/nls/example/nls"
	"golang.org/x/text/language"
)

//go:generate nls -dir messages -pkg nls -v
func main() {
	loc := New(language.Dutch.String(), language.English.String())

	fmt.Println(loc.Get(M_hello))
	fmt.Println(loc.Get(M_world))
	fmt.Println(loc.Get(M_sky, "Sky"))
	fmt.Println(loc.Get("ying", "yang"))
	fmt.Println(loc.Format(M_sea1, "name", "Noord"))
	fmt.Println(loc.Format(M_cats1, "count", 3))
	fmt.Println(loc.Format(M_cats1, "count", 1))
}
