package scanner

import (
	"slices"
	"testing"
)

func Test_ParseScanemas(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		input            string
		expectedScanemas []Scanema
		expectedUnparsed []byte
	}{
		{"Basic Case 1", "Hello, world! ", []Scanema{
			{Content: []byte("Hello,"), ScanemaType: ScanemaTypeCommonWord},
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
			{Content: []byte("world!"), ScanemaType: ScanemaTypeCommonWord},
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
		}, nil},
		{"EOL Case 1", "\n", []Scanema{
			{Content: []byte("\n"), ScanemaType: ScanemaTypeEOL},
		}, nil},
		{"EOL Case 2", "\r\n", []Scanema{
			{Content: []byte("\r\n"), ScanemaType: ScanemaTypeEOL},
		}, nil},
		{"Tab Case", "\t", []Scanema{
			{Content: []byte("\t"), ScanemaType: ScanemaTypeSpecialSymbol},
		}, nil},
		{"Spaces Case", "  ", []Scanema{
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
		}, nil},
		{"Carriage Return Case", "\r ", []Scanema{
			{Content: []byte("\r"), ScanemaType: ScanemaTypeCommonWord},
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
		}, nil},
		{"Tab + EOL Case", "\t\n", []Scanema{
			{Content: []byte("\t"), ScanemaType: ScanemaTypeSpecialSymbol},
			{Content: []byte("\n"), ScanemaType: ScanemaTypeEOL},
		}, nil},
		{"Carriage Return + Tab + EOL Case", "\r\n\t\n", []Scanema{
			{Content: []byte("\r\n"), ScanemaType: ScanemaTypeEOL},
			{Content: []byte("\t"), ScanemaType: ScanemaTypeSpecialSymbol},
			{Content: []byte("\n"), ScanemaType: ScanemaTypeEOL},
		}, nil},
		// Edge cases
		{"Empty Input", "", nil, nil},
		{"EOL Case 3", "\r\n", []Scanema{
			{Content: []byte("\r\n"), ScanemaType: ScanemaTypeEOL},
		}, nil},
		// Multi-byte characters
		{"Japanese Characters", "こんにちは ", []Scanema{
			{Content: []byte("こんにちは"), ScanemaType: ScanemaTypeCommonWord},
			{Content: []byte(" "), ScanemaType: ScanemaTypeSpecialSymbol},
		}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			scanemas, unparsed := parseScanemas([]byte(tc.input))
			for i, sc := range tc.expectedScanemas {
				if i >= len(scanemas) {
					t.Errorf("Expected scanemas: %v, got: %v", tc.expectedScanemas, scanemas)
					return
				}
				if string(scanemas[i].Content) != string(sc.Content) || scanemas[i].ScanemaType != sc.ScanemaType {
					t.Errorf("Expected scanemas: %v, got: %v", tc.expectedScanemas, scanemas)
				}
			}
			if !slices.Equal(unparsed, tc.expectedUnparsed) {
				t.Errorf("Expected unparsed bytes: %v, got: %v", tc.expectedUnparsed, unparsed)
			}
		})
	}
}

func Test_GetScanemaType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		word         string
		expectedType ScanemaType
	}{
		{"Empty String", "", ScanemaTypeEmpty},
		{"Common Word", "Hello", ScanemaTypeCommonWord},
		{"Newline", "\n", ScanemaTypeEOL},
		{"Tab", "\t", ScanemaTypeCommonWord},
		{"Carriage Return", "\r", ScanemaTypeCommonWord},
		{"Carriage Return + Newline", "\r\n", ScanemaTypeEOL},
		{"Tab + Newline", "\t", ScanemaTypeCommonWord},
		{"Carriage Return + Newline + Tab + Newline", "\n", ScanemaTypeEOL},
		{"Japanese Characters", "こんにちは", ScanemaTypeCommonWord},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actualType := getScanemaType(tc.word)
			if actualType != tc.expectedType {
				t.Errorf("Expected ScanemaType for \"%s\": %d, got: %d", tc.word, tc.expectedType, actualType)
			}
		})
	}
}

func Test_In(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		toFind        []byte
		searchTargets []SpecialSymbol
		expected      bool
	}{
		{"Not in Whitespace", []byte{'H', 'e', 'l', 'l', 'o'}, []SpecialSymbol{Tab, Space, NewLine}, false},
		{"Whitespace", []byte{' '}, []SpecialSymbol{Tab, Space, NewLine}, true},
		{"Newline", []byte{'\n'}, []SpecialSymbol{Tab, Space, NewLine}, true},
		{"Carriage Return & Newline", []byte{'\r', '\n'}, []SpecialSymbol{NewLine, WinNewLine}, true},
		{"Tab", []byte{'\t'}, []SpecialSymbol{Tab, Space, NewLine}, true},
		{"Not in Single Carriage Return", []byte{'\r'}, []SpecialSymbol{Tab, Space, NewLine}, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := in(tc.toFind, tc.searchTargets)
			if actual != tc.expected {
				t.Errorf("Expected in result for %v and %v: %v, got: %v", tc.toFind, tc.searchTargets, tc.expected, actual)
			}
		})
	}
}
