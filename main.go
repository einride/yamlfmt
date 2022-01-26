package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", "directory to format")
	var recursive bool
	flag.BoolVar(&recursive, "r", false, "recursive flag for directory")
	var file string
	flag.StringVar(&file, "file", "", "file to format")
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()
	if len(os.Args[1:]) == 0 {
		flag.Usage()
	}
	if dir != "" && file != "" {
		panic("Pick dir or file, not both")
	}
	if err := format(dir, file, recursive, verbose); err != nil {
		panic(err)
	}
}

func format(dir, file string, recursive, verbose bool) error {
	if file != "" {
		return formatFile(file, verbose)
	}
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			if recursive {
				if err := formatFile(path, verbose); err != nil {
					return err
				}
			} else {
				if filepath.Dir(path) == dir {
					if err := formatFile(path, verbose); err != nil {
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

func formatFile(path string, verbose bool) error {
	if verbose {
		log.Printf("Formatting %s", path)
	}
	node := yaml.Node{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	yamlFile = preserveEmptyLines(yamlFile)
	if err := yaml.Unmarshal(yamlFile, &node); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	if err := encoder.Encode(&node); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return ioutil.WriteFile(path, cleanupPreserveEmptyLines(b.Bytes()), 0o600)
}
