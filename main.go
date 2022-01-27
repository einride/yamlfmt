package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var dir string
	flag.StringVar(&dir, "d", "", "directory to formatRecursive")
	var recursive bool
	flag.BoolVar(&recursive, "r", false, "recursive flag for directory")
	var file string
	flag.StringVar(&file, "f", "", "file to formatRecursive")
	var indent int
	flag.IntVar(&indent, "i", 2, "indentation")
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()
	if len(os.Args[1:]) == 0 {
		flag.Usage()
	}
	if dir != "" && file != "" {
		log.Fatalln("Pick dir or file, not both")
	}
	if err := formatRecursive(dir, file, indent, recursive, verbose); err != nil {
		log.Fatalln(err)
	}
}

func formatRecursive(dir, file string, indent int, recursive, verbose bool) error {
	formatFile := func(path string) error {
		if verbose {
			log.Printf("Formatting %s", path)
		}
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b, err := format(yamlFile, indent)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		return os.WriteFile(path, b, 0o600)
	}

	if file != "" {
		return formatFile(file)
	}

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			if recursive {
				if err := formatFile(path); err != nil {
					return err
				}
			} else {
				if filepath.Dir(path) == dir {
					if err := formatFile(path); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
