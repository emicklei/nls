package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

var (
	oDir     = flag.String("dir", "", "directory to scan for .yaml files")
	oPkg     = flag.String("pkg", "nls", "package name for the generated code")
	oVerbose = flag.Bool("v", false, "verbose output")
)

// go run . -v -dir ../../example/messages -pkg ../../example/nls
func main() {
	flag.Parse()
	langDirs, err := os.ReadDir(*oDir)
	if err != nil {
		log.Fatal(err)
	}
	allEntries := []Entry{}
	for _, each := range langDirs {
		if each.IsDir() {
			messageFiles, err := os.ReadDir(filepath.Join(*oDir, each.Name()))
			if err != nil {
				log.Printf("cannot read directory %s\n", each.Name())
			}
			for _, file := range messageFiles {
				if filepath.Ext(file.Name()) == ".yaml" {
					fullName := filepath.Join(*oDir, each.Name(), file.Name())
					if entries, err := collectEntries(each.Name(), fullName); err != nil {
						log.Printf("cannot process file %s:%v\n", err, fullName)
					} else {
						allEntries = append(allEntries, entries...)
					}
				}
			}
		}
	}
	if *oVerbose {
		for _, each := range allEntries {
			log.Printf("%s.%s=%s\n", each.Language, each.Key, each.Text)
		}
	}
	if err := os.Mkdir(*oPkg, os.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatalf("%[1]T %[1]v", err)
	}
	if err := writeGoFile(allEntries); err != nil {
		log.Fatal(err)
	}
	if err := writeEntries(allEntries, *oDir); err != nil {
		log.Fatal(err)
	}
}

//go:embed localizer.template
var localizerTemplate string

// constantName returns the full constant name, e.g. M_hello, M_cats1, M_trends2_3
func constantName(e Entry) string {
	name := "M_" + e.Key
	// count replacements
	replacements := strings.Count(e.Text, "{{.")
	if replacements == 0 {
		return name
	}
	// check if last rune of key is a digit
	runes := []rune(e.Key)
	if len(runes) > 0 && unicode.IsDigit(runes[len(runes)-1]) {
		return fmt.Sprintf("%s_%d", name, replacements)
	}
	return fmt.Sprintf("%s%d", name, replacements)
}

func writeGoFile(entries []Entry) error {
	outName := filepath.Join(*oPkg, "generated_catalog.go")
	if *oVerbose {
		log.Printf("writing %s\n", outName)
	}
	out, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer out.Close()
	tmpl, err := template.New("localizer").Funcs(template.FuncMap{
		"constantName": constantName,
	}).Parse(localizerTemplate)
	if err != nil {
		return err
	}
	uniqueEntries := map[string]Entry{}
	// collect unique entries
	for _, each := range entries {
		if _, ok := uniqueEntries[each.Key]; !ok {
			uniqueEntries[each.Key] = each
		}
	}
	data := struct {
		Package       string
		UniqueEntries map[string]Entry
		Entries       []Entry
	}{
		Package:       filepath.Base(*oPkg),
		UniqueEntries: uniqueEntries,
		Entries:       entries,
	}
	return tmpl.Execute(out, data)
}

func collectEntries(language, fullName string) ([]Entry, error) {
	if *oVerbose {
		log.Printf("processing %s in [%s]\n", fullName, language)
	}
	reader, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	dec := yaml.NewDecoder(reader)
	messages := make(map[string]any)
	err = dec.Decode(&messages)
	if err != nil {
		return nil, err
	}
	if *oVerbose {
		log.Printf("%d messages found\n", len(messages))
	}
	var entries []Entry
	for key, value := range messages {
		if s, ok := value.(string); ok {
			entries = append(entries, Entry{Language: language, Key: key, Text: s})
		} else if m, ok := value.(map[string]any); ok {
			entry := Entry{Language: language, Key: key}
			if v, ok := m["value"].(string); ok {
				entry.Text = v
			}
			if v, ok := m["description"].(string); ok {
				entry.Description = v
			}
			entries = append(entries, entry)
		}
	}
	return entries, nil
}
