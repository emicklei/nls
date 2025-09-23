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
	"slices"
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
	allEntries = fillMissingEntries(allEntries)
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
	replacements := e.Replacements()
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
	languages := []string{}
	// collect unique entries
	for _, each := range entries {
		// TODO use map?
		if !slices.Contains(languages, each.Language) {
			languages = append(languages, each.Language)
		}
		existing, ok := uniqueEntries[each.Key]
		if !ok {
			uniqueEntries[each.Key] = each
		} else {
			// give priority to entry with description
			if len(each.Description) > 0 && len(existing.Description) == 0 {
				uniqueEntries[each.Key] = each
			}
		}
	}
	data := struct {
		Package       string
		UniqueEntries map[string]Entry
		Entries       []Entry
		LanguageTags  []string
	}{
		Package:       filepath.Base(*oPkg),
		UniqueEntries: uniqueEntries,
		Entries:       entries,
		LanguageTags:  languages,
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
	var node yaml.Node
	err = dec.Decode(&node)
	if err != nil {
		return nil, err
	}
	if *oVerbose {
		log.Printf("%d messages found\n", len(node.Content))
	}
	var entries []Entry
	for i, content := range node.Content[0].Content {
		// is key?
		if i%2 == 0 {
			keyNode := content
			valueNode := node.Content[0].Content[i+1]
			if valueNode.Tag == "!!str" {
				entries = append(entries, Entry{Language: language, Key: keyNode.Value, Text: valueNode.Value, Comment: keyNode.HeadComment})
			} else if valueNode.Tag == "!!map" {
				entry := Entry{Language: language, Key: keyNode.Value, Comment: keyNode.HeadComment}
				for j, each := range valueNode.Content {
					if j%2 == 0 {
						mapkeyNode := each
						mapvalueNode := valueNode.Content[j+1]
						if mapkeyNode.Value == "msg" {
							entry.Text = mapvalueNode.Value
						}
						if mapkeyNode.Value == "desc" {
							entry.Description = mapvalueNode.Value
						}
					}
				}
				entries = append(entries, entry)
			}
		}
	}
	return entries, nil
}

func fillMissingEntries(allEntries []Entry) []Entry {
	// key is language
	entriesPerLanguage := map[string][]Entry{}
	for _, each := range allEntries {
		entriesPerLanguage[each.Language] = append(entriesPerLanguage[each.Language], each)
	}
	// key is message key
	allKeys := map[string]Entry{}
	for _, each := range allEntries {
		// only set description if not already set or empty
		if d, ok := allKeys[each.Key]; !ok || d.Description == "" {
			allKeys[each.Key] = each
		}
	}
	for lang, entries := range entriesPerLanguage {
		// key is message key
		keysInLang := map[string]bool{}
		for _, each := range entries {
			keysInLang[each.Key] = true
		}
		for key, entryWithInfo := range allKeys {
			if !keysInLang[key] {
				// missing key
				if *oVerbose {
					log.Printf("language [%s] is missing key [%s]", lang, key)
				}
				newEntry := Entry{
					Language:    lang,
					Key:         key,
					Description: entryWithInfo.Description,
					Comment:     entryWithInfo.Comment,
				}
				allEntries = append(allEntries, newEntry)
			}
		}
	}
	return allEntries
}
