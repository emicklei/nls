package main

import (
	_ "embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"text/template"

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

func writeGoFile(entries []Entry) error {
	outName := filepath.Join(*oPkg, "localizer.go")
	if *oVerbose {
		log.Printf("writing %s\n", outName)
	}
	out, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer out.Close()
	tmpl, err := template.New("localizer").Parse(localizerTemplate)
	if err != nil {
		return err
	}
	keys := []string{}
	// collect unique keys
	for _, each := range entries {
		if slices.Contains(keys, each.Key) {
			continue
		}
		keys = append(keys, each.Key)
	}
	data := struct {
		Package string
		Keys    []string
		Entries []Entry
	}{
		Package: filepath.Base(*oPkg),
		Keys:    keys,
		Entries: entries,
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
	messages := make(map[string]string)
	err = dec.Decode(&messages)
	if err != nil {
		return nil, err
	}
	if *oVerbose {
		log.Printf("%d messages found\n", len(messages))
	}
	var entries []Entry
	for key, value := range messages {
		entries = append(entries, Entry{Language: language, Key: key, Text: value})
	}
	return entries, nil
}
