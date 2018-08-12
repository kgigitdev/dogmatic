package godgt

import (
	"fmt"
	"strings"
)

// FenBuilder builds up a "mini FEN" (that is, just the positional
// part of the FEN, not the metadata). It is expected that it be
// called precisely 64 times, once for each square in the board, in
// the order a8...h8, a7...h7, [...], a1...h1. Non-empty squares
// should be called with the relevant piece code (KQRBNPkqrbnp).
// Empty squares should be called with an empty string or an
// underscore. After the 64th invocation, String() will return the
// FEN.
//
// Note that this FEN dialect doesn't elide empty ranks, nor does it
// elide trailing digits. For example, "/8/" could also be written
// "//", and "/4p3/" could also be written as "/4p/". This is by
// design, as we aim to match notnil/chess's FEN dialect (to match
// positions).
type FenBuilder struct {
	elems []string
	file  int
	rank  int
	fen   string
}

// NewFenBuilder creates and returns a pointer to a new FenBuilder instance.
func NewFenBuilder() *FenBuilder {
	return &FenBuilder{
		elems: make([]string, 0),
		file:  1,
		rank:  8,
	}
}

func (fb *FenBuilder) String() string {
	fb.collapseEmptySpaces()
	return fb.fen
}

// Add accepts a single square's character. It needs to be called
// exactly 64 times.
func (fb *FenBuilder) Add(fenChar string) {
	if fb.rank < 1 {
		return
	}
	if fenChar == "" {
		fenChar = "_"
	}
	fb.elems = append(fb.elems, fenChar)
	fb.file++
	if fb.file > 8 {
		if fb.rank > 1 {
			fb.elems = append(fb.elems, "/")
		}
		fb.file = 1
		fb.rank--
	}
}

func (fb *FenBuilder) collapseEmptySpaces() {
	spaces := 0
	collapsed := make([]string, 0)
	for _, elem := range fb.elems {
		if elem != "_" {
			if spaces > 0 {
				d := fmt.Sprintf("%d", spaces)
				collapsed = append(collapsed, d)
				spaces = 0
			}
			collapsed = append(collapsed, elem)
		} else {
			spaces++
		}
	}

	if spaces > 0 {
		d := fmt.Sprintf("%d", spaces)
		collapsed = append(collapsed, d)
		spaces = 0
	}

	fb.fen = strings.Join(collapsed, "")
}
