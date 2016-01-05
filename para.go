// Wrap text to a given column number. Acts as a filter from stdin to stdout.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const defaultWrapCol = 80

func main() {
	var wrap int
	if len(os.Args) < 2 {
		wrap = defaultWrapCol
	} else {
		num, err := strconv.ParseInt(os.Args[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		} else {
			wrap = int(num)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	r := Rapper{maxcols: wrap}
	err := r.Wraptext(scanner, writer)
	if err != nil {
		log.Fatal(err)
	}
}

// Rapper contains methods for compressive text wrapping
type Rapper struct {
	maxcols      int
	carry        int
	pendingBreak bool
}

// wrap text to column length, compact paragraph along the way
// respect lines that end in a period
func (r Rapper) Wraptext(scanner *bufio.Scanner, writer *bufio.Writer) error {
	for scanner.Scan() {
		if r.pendingBreak {
			writer.WriteString("\n")
			r.pendingBreak = false
		}
		line := scanner.Text()
		wrapped := r.wrapline(line)
		writer.WriteString(wrapped)
		if strings.HasSuffix(line, ".") || len(line) == 0 {
			// Respect paragraphs and full stops
			r.pendingBreak = true
			r.carry = 0
		} else if strings.HasSuffix(wrapped, ".") {
			r.carry = 0
		} else {
			lastbrk := strings.LastIndex(wrapped, "\n")
			r.carry = len(line) - 1 - lastbrk
			if r.carry > 0 {
				r.carry = r.carry + 1
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return writer.Flush()
}

// wrap a single line to a colum length, possibly breaking it
func (r Rapper) wrapline(line string) string {
	lastWhite, lastNewline := -1, -r.carry-1
	out := make([]byte, len(line))
	var initbreak bool
	if len(line) == 0 && r.carry > 0 {
		initbreak = true
	}
	for j := 0; j < len(line); j++ {
		out[j] = line[j]
		if unicode.IsSpace(rune(line[j])) {
			lastWhite = j
		}
		if j-lastNewline > r.maxcols && lastWhite > -1 {
			out[lastWhite] = '\n'
			lastNewline = lastWhite
		} else if j-lastNewline > r.maxcols {
			initbreak = true
			lastNewline = lastWhite
		}
	}
	if initbreak {
		return "\n" + string(out)
	} else if r.carry > 0 {
		return " " + string(out)
	} else {
		return string(out)
	}
}
