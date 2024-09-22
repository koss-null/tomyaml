package scanner

type ScanemaType byte

const (
	ScanemaTypeEmpty         ScanemaType = 0
	ScanemaTypeCommonWord    ScanemaType = 1
	ScanemaTypeSpecialSymbol ScanemaType = 2
	ScanemaTypeEOL           ScanemaType = 3
)

// Scanema is a valid utf-8 minimal string token
type Scanema struct {
	Content []byte
	ScanemaType
}

type SpecialSymbol []byte

var (
	Tab        SpecialSymbol = []byte{'\t'}
	Space      SpecialSymbol = []byte{' '}
	NewLine    SpecialSymbol = []byte{'\n'}
	WinNewLine SpecialSymbol = []byte{'\r', '\n'}
)

var (
	whitespaceSymbols          = []SpecialSymbol{Tab, Space}
	newLineSymbols             = []SpecialSymbol{NewLine}
	doubleSymbolNewLineSymbols = []SpecialSymbol{WinNewLine}
)

// parseScanemas returns the slice of parsed scanemas and a slice of bytes that are not parsed
// (so the end of the scanema is expected in the next lines)
func parseScanemas(lines []byte) ([]Scanema, []byte) {
	scnStart, scnEnd := 0, 0
	var prevByte *byte
	// FIXME: 64
	scanemas := make([]Scanema, 0, 64)
	for {
		// last scanema in lines
		if scnEnd >= len(lines) {
			return scanemas, lines[scnStart:]
		}

		curByte := lines[scnEnd]
		scnEnd++

		switch {
		case in([]byte{curByte}, whitespaceSymbols):
			if scnEnd-scnStart > 1 {
				// append prev scanema as a common word
				scanemas = append(scanemas, Scanema{
					Content:     lines[scnStart : scnEnd-1],
					ScanemaType: ScanemaTypeCommonWord,
				})
			}
			// append current symbol as a space scanema
			scanemas = append(scanemas, Scanema{
				Content:     []byte(string(curByte)),
				ScanemaType: ScanemaTypeSpecialSymbol,
			})
			scnStart = scnEnd
			// multiple bytes new line
		case prevByte != nil && in([]byte{*prevByte, curByte}, doubleSymbolNewLineSymbols):
			// append prev scanema as common word
			if scnEnd-scnStart > 2 {
				scanemas = append(scanemas, Scanema{
					Content:     lines[scnStart : scnEnd-2],
					ScanemaType: ScanemaTypeCommonWord,
				})
			}
			// append current symbol
			scanemas = append(scanemas, Scanema{
				Content:     lines[scnEnd-2 : scnEnd],
				ScanemaType: ScanemaTypeEOL,
			})
			scnStart = scnEnd

		// single byte new line
		case in([]byte{curByte}, newLineSymbols):
			if scnEnd-scnStart > 1 {
				// append prev scanema as common word
				scanemas = append(scanemas, Scanema{
					Content:     lines[scnStart : scnEnd-1],
					ScanemaType: ScanemaTypeCommonWord,
				})
			}
			// append current symbol
			scanemas = append(scanemas, Scanema{
				Content:     []byte(string(curByte)),
				ScanemaType: ScanemaTypeEOL,
			})
			scnStart = scnEnd
		}
		prevByte = &curByte
	}
}

func getScanemaType(word string) ScanemaType {
	if len(word) == 0 {
		return ScanemaTypeEmpty
	}

	runeWord := []byte(word)
	switch {
	case in(runeWord, whitespaceSymbols):
		return ScanemaTypeCommonWord
	case in(runeWord, newLineSymbols), in(runeWord, doubleSymbolNewLineSymbols):
		return ScanemaTypeEOL
	default:
		return ScanemaTypeCommonWord
	}
}

func in(toFind []byte, searchTargets []SpecialSymbol) bool {
	for _, st := range searchTargets {
		found := true
		for i := range toFind {
			if len(st) != len(toFind) || toFind[i] != st[i] {
				found = false
				break
			}
		}
		if found {
			return true
		}
	}
	return false
}
