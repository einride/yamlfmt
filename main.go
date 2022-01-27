package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
		if verbose {
			log.Printf("Formatting %s", path)
		}
		node := yaml.Node{}
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		yamlFile = preserveEmptyLines(yamlFile)
		if err := yaml.Unmarshal(yamlFile, &node); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		var b bytes.Buffer
		encoder := yaml.NewEncoder(&b)
		encoder.SetIndent(indent)
		if err := encoder.Encode(&node); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		return os.WriteFile(path, cleanupPreserveEmptyLines(b.Bytes()), 0o600)
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

// preserveEmptyLines adds a temporary #comment on each empty line in the provided byte array.
// cleanupPreserveEmptyLines can be used to clean up the temporary comments.
func preserveEmptyLines(src []byte) []byte {
	return bytes.ReplaceAll(src, []byte("\n\n"), []byte("\n#preserveEmptyLine\n"))
}

// cleanupPreserveEmptyLines cleans up the temporary #comment added by PreserveEmptyLines.
func cleanupPreserveEmptyLines(src []byte) []byte {
	// Remove temporary comment.
	indentPreserveComment := regexp.MustCompile("\n\\s+#preserveEmptyLine\n")
	src = indentPreserveComment.ReplaceAll(src, []byte("\n\n"))
	src = bytes.ReplaceAll(src, []byte("\n#preserveEmptyLine\n"), []byte("\n\n"))
	// Remove trailing empty lines
	src = bytes.TrimSpace(src)
	src = append(src, []byte("\n")...)
	return src
}
