package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	var dir string
	flag.StringVar(&dir, "d", "", "directory to format")
	var recursive bool
	flag.BoolVar(&recursive, "r", false, "recursive flag for directory")
	var file string
	flag.StringVar(&file, "f", "", "file to format")
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
	if err := format(dir, file, indent, recursive, verbose); err != nil {
		log.Fatalln(err)
	}
}

func format(dir, file string, indent int, recursive, verbose bool) error {
	formatFile := func(path string) error {
		var node yaml.Node
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		commentedFile := preserveEmptyLines(yamlFile)
		if err := yaml.Unmarshal(commentedFile, &node); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		var b bytes.Buffer
		encoder := yaml.NewEncoder(&b)
		encoder.SetIndent(indent)
		if err := encoder.Encode(&node); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if equal := bytes.Compare(yamlFile, cleanupPreserveEmptyLines(b.Bytes())); equal != 0 {
			if verbose {
				log.Printf("Formatted %s", path)
			}
			return os.WriteFile(path, cleanupPreserveEmptyLines(b.Bytes()), 0o600)
		}
		return nil
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
				if filepath.Dir(path) == filepath.Clean(dir) {
					if err := formatFile(path); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

// preserveEmptyLines adds a temporary #comment on each empty line in the provided byte array.
// cleanupPreserveEmptyLines can be used to clean up the temporary comments.
func preserveEmptyLines(src []byte) []byte {
	var b []byte
	in := bufio.NewScanner(bytes.NewReader(src))
	for in.Scan() {
		line := in.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			b = append(b, []byte("#preserveEmptyLine\n")...)
		} else {
			b = append(b, append(line, "\n"...)...)
		}
	}
	return b
}

// cleanupPreserveEmptyLines cleans up the temporary #comment added by PreserveEmptyLines.
func cleanupPreserveEmptyLines(src []byte) []byte {
	// Scrub one time to preserve multi line strings
	var first []byte
	in := bufio.NewScanner(bytes.NewReader(src))
	for in.Scan() {
		line := in.Text()
		if strings.TrimSpace(line) == "#preserveEmptyLine" {
			first = append(first, []byte("\n")...)
			continue
		}
		if strings.Contains(line, "#preserveEmptyLine") {
			line = strings.ReplaceAll(line, "#preserveEmptyLine", "\n#preserveEmptyLine")
			first = append(first, []byte(line+"\n")...)
			continue
		}
		first = append(first, []byte(line+"\n")...)
	}
	// Scrub a second time to remove the rest
	var second []byte
	in = bufio.NewScanner(bytes.NewReader(first))
	for in.Scan() {
		line := in.Text()
		if strings.TrimSpace(line) == "#preserveEmptyLine" {
			second = append(second, []byte("\n")...)
			continue
		}
		if strings.Contains(line, "#preserveEmptyLine") {
			line = strings.ReplaceAll(line, "#preserveEmptyLine", "\n")
			second = append(second, []byte(line+"\n")...)
			continue
		}
		second = append(second, []byte(line+"\n")...)
	}
	// Remove trailing empty lines
	second = bytes.TrimSpace(second)
	second = append(second, []byte("\n")...)
	return second
}
