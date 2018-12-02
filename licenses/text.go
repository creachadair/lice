package licenses

import (
	"strings"
	"unicode"
)

func newBlock(s string) *block {
	return &block{lines: strings.Split(s, "\n")}
}

type block struct {
	lines []string
}

func (t *block) trimSpace() *block {
	// Discard trailing whitespace from all lines.
	for i, line := range t.lines {
		t.lines[i] = strings.TrimRight(line, " \t\r\n")
	}

	// Discard leading and trailing blank lines.
	i := 0
	for i < len(t.lines) && t.lines[i] == "" {
		i++
	}
	j := len(t.lines)
	for j > i && t.lines[j-1] == "" {
		j--
	}
	t.lines = t.lines[i:j]
	return t
}

func (t *block) untabify(width int) *block {
	if width <= 0 {
		width = 4
	}
	spaces := strings.Repeat(" ", width)
	for i, line := range t.lines {
		t.lines[i] = strings.Replace(line, "\t", spaces, -1)
	}
	return t
}

func (t *block) leftJust() *block {
	var min *string
	for _, line := range t.lines {
		spc := leftSpace(line)
		if line != "" && (min == nil || len(spc) < len(*min)) {
			min = &spc
		}
	}
	if min != nil {
		for i, line := range t.lines {
			t.lines[i] = strings.TrimPrefix(line, *min)
		}
	}
	return t
}

func (t *block) indent(ind string) *block {
	for i, line := range t.lines {
		if line == "" && i+1 == len(t.lines) {
			continue
		}
		t.lines[i] = strings.TrimRight(ind+line, " \t\r\n")
	}
	return t
}

func (t *block) prepend(s string) *block {
	t.lines = append([]string{s}, t.lines...)
	return t
}

func (t *block) append(ss ...string) *block {
	t.lines = append(t.lines, ss...)
	return t
}

func (t *block) String() string { return strings.Join(t.lines, "\n") }

func leftSpace(s string) string {
	var left string
	for _, c := range s {
		if !unicode.IsSpace(c) {
			break
		}
		left += string(c)
	}
	return left
}

// An Indenting is a rule for indenting or commenting license text for
// insertion into a file. A nil Indenting leaves the input text unmodified.
type Indenting func(*block) *block

func (in Indenting) fix(b *block) *block {
	if in == nil {
		return b
	}
	return in(b)
}

// IPrefix constructs an Indenting that prefixes each line of text with the
// specified marker.
func IPrefix(marker string) Indenting {
	return func(b *block) *block { return b.indent(marker) }
}

// IComment constructs an Indenting that prefixes the lines of text with the
// given comment markers.
func IComment(first, rest, last string) Indenting {
	return func(b *block) *block {
		return b.indent(rest).prepend(first).append(last)
	}
}
