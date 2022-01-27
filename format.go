package main

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

func format(src []byte, indent int) ([]byte, error) {
	if len(src) == 0 {
		return src, nil
	}
	src = preserveEmptyLines(src)
	node := yaml.Node{}
	if err := yaml.Unmarshal(src, &node); err != nil {
		return nil, fmt.Errorf("%s: %w", string(src), err)
	}
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(indent)
	if err := encoder.Encode(&node); err != nil {
		return nil, err
	}
	return cleanupPreserveEmptyLines(b.Bytes()), nil
}

// preserveEmptyLines adds a temporary #comment on each empty line in the provided byte array.
// cleanupPreserveEmptyLines can be used to clean up the temporary comments.
func preserveEmptyLines(src []byte) []byte {
	return bytes.ReplaceAll(src, []byte("\n\n"), []byte("\n#preserveEmptyLine\n"))
}

// cleanupPreserveEmptyLines cleans up the temporary #comment added by PreserveEmptyLines.
func cleanupPreserveEmptyLines(src []byte) []byte {
	x := string(src)
	_ = x
	// Remove temporary comment.
	indentPreserveComment := regexp.MustCompile("\n\\s+#preserveEmptyLine\n")
	src = indentPreserveComment.ReplaceAll(src, []byte("\n\n"))
	src = bytes.ReplaceAll(src, []byte("\n#preserveEmptyLine\n"), []byte("\n\n"))
	// Remove trailing empty lines
	src = bytes.TrimSpace(src)
	src = append(src, []byte("\n")...)
	dst := string(src)
	_ = dst
	return src
}
