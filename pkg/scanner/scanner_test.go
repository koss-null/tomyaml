package scanner

import (
	"slices"
	"strings"
	"testing"
)

func TestScan_ValidTOML(t *testing.T) {
	// Sample TOML content
	tomlContent := `
		[owner]
		name = "Tom Preston-Werner"
		dob = 1979-05-27T07:32:00Z
		`

	// Create a reader from the TOML content
	reader := strings.NewReader(tomlContent)

	// Call the Scan function
	scanemas, err := Scan(reader)
	// Assert no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that we received a non-nil result
	if scanemas == nil {
		t.Errorf("Expected non-nil scanemas")
	}

	scns := make([]Scanema, 0)
	for scn, ok := scanemas.PopFront(); ok; scn, ok = scanemas.PopFront() {
		scns = append(scns, scn...)
	}

	expectedScanemas := []Scanema{
		{
			Content:     []byte{'\n'},
			ScanemaType: ScanemaTypeEOL,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("[owner]"),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{'\n'},
			ScanemaType: ScanemaTypeEOL,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("name"),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{' '},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("="),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{' '},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("\"Tom"),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{' '},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("Preston-Werner\""),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{'\n'},
			ScanemaType: ScanemaTypeEOL,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("dob"),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{' '},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("="),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{' '},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte("1979-05-27T07:32:00Z"),
			ScanemaType: ScanemaTypeCommonWord,
		},
		{
			Content:     []byte{'\n'},
			ScanemaType: ScanemaTypeEOL,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
		{
			Content:     []byte{'\t'},
			ScanemaType: ScanemaTypeSpecialSymbol,
		},
	}

	// Check that each Scanema in scns matches the expected values
	for i, scn := range scns {
		if len(expectedScanemas) <= i {
			t.Errorf("not expected: %v", scns[i:])
			return
		}
		if !slices.Equal(scn.Content, expectedScanemas[i].Content) &&
			scn.ScanemaType == expectedScanemas[i].ScanemaType {
			t.Errorf(
				"Scanema at index %d does not match expected value: %q expected: %q",
				i,
				string(scn.Content),
				string(expectedScanemas[i].Content),
			)
		}
	}
}
