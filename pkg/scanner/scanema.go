package scanner

import (
	"strings"
	"unicode/utf8"
)

type ScanemaType int

const (
	ScanemaTypeCommonWord    ScanemaType = 1
	ScanemaTypeSpecialSymbol ScanemaType = 2
	ScanemaTypeEOLWord       ScanemaType = 3
)

type Scanema struct {
	Content string
	ScanemaType
}

type SpecialSymbol string

const (
	Tab        SpecialSymbol = "\t"
	Space      SpecialSymbol = " "
	NewLine    SpecialSymbol = "\n"
	WinNewLine SpecialSymbol = "\r\n"
)

var (
	whitespaceSymbols = []SpecialSymbol{Tab, Space}
	newLineSymbols    = []SpecialSymbol{NewLine, WinNewLine}
)

func parseScanemas(lines string) []Scanema {
	scnStart, scnEnd := 0, 0
	scanemas := make([]Scanema, 1)
	for {
		// last word in lines
		if len(lines) >= scnEnd {
			word := lines[scnStart:]
			return append(scanemas, Scanema{
				Content:     word,
				ScanemaType: getScanemaType(word),
			})
		}

		runeVal, width := utf8.DecodeRuneInString(lines[scnEnd:])
		scnEnd += width

		switch {
		case in(runeVal, whitespaceSymbols):
			scanemas = append(scanemas, Scanema{
				Content:     lines[scnStart:scnEnd],
				ScanemaType: ScanemaTypeCommonWord,
			})
			scanemas = append(scanemas, Scanema{
				Content:     string(runeVal),
				ScanemaType: ScanemaTypeSpecialSymbol,
			})
			scnStart = scnEnd

		case in(runeVal, newLineSymbols):
			scanemas = append(scanemas, Scanema{
				Content:     lines[scnStart:scnEnd],
				ScanemaType: ScanemaTypeEOLWord,
			})
			scanemas = append(scanemas, Scanema{
				Content:     string(NewLine),
				ScanemaType: ScanemaTypeSpecialSymbol,
			})
			scnStart = scnEnd
		}
	}
}

func getScanemaType(word string) ScanemaType {
	if len(word) == 0 {
		return ScanemaTypeCommonWord
	}

	// this step is not necessary for now, but may be useful to
	// escape some mistakes to be made on this func change
	if in(rune(word[len(word)-1]), whitespaceSymbols) {
		return ScanemaTypeCommonWord
	}
	if strings.HasSuffix(word, string(NewLine)) || strings.HasSuffix(word, string(WinNewLine)) {
		return ScanemaTypeEOLWord
	}
	return ScanemaTypeCommonWord
}

func in(toFind rune, searchTargets []SpecialSymbol) bool {
	for i := range searchTargets {
		if strings.ContainsRune(string(searchTargets[i]), toFind) {
			return true
		}
	}
	return false
}
