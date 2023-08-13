package imapgoprebuilt

import (
	"bufio"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func DecodeLines(html string) ([]string, error) {
	file := strings.NewReader(html)

	decodingReader := transform.NewReader(file, charmap.Windows1252.NewDecoder())

	lines := []string{}

	scanner := bufio.NewScanner(decodingReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
