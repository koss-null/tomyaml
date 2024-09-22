# Scanner Package 
  
## Purpose:
This package is designed to efficiently **scan text files** and **identify various lexemes** (meaningful units) within the text.
  
## Key Features:
  
- Identifies common words, special symbols, and end-of-line characters (EOL).
- Processes text files in chunks for optimized memory usage.
  
## Usage:
  
The primary function for scanning a file is `Scan(file io.Reader)`.
  
The `Scan` function returns two values:
  
- A pointer to a `list.Linked` containing slices of Scanema objects.
Each slice represents a chunk of scanned data from the file.
- An `error` object, if any errors occur during the scanning process.
  
## Scanema Types:
  
The **Scanema** struct represents a **single unit** identified during the scan. It has two fields:
  
- `Content`: A byte slice containing the actual content of the scan unit.
- `ScanemaType`: An enum value indicating the unit's type.
  
Valid types are:
- *ScanemaTypeEmpty*: An empty unit (no content).
- *ScanemaTypeCommonWord*: A sequence of characters that is not a special symbol or EOL.
- *ScanemaTypeSpecialSymbol*: A single special symbol like whitespace (`\t` or ` `).
- *ScanemaTypeEOL*: An end-of-line character, either `\n` (Unix) or `\r\n` (Windows).
  
## Corner Cases and Considerations:
  
**Multiple Whitespace Symbols**: Consecutive whitespace characters are treated as separate ScanemaTypeSpecialSymbol instances.  
**Windows Newlines**: The scanner specifically handles Windows newline format (\r\n) as a single ScanemaTypeEOL.  
**Partial Data at End of File**: To ensure the last word of a file is parsed correctly, a \r\n sequence is appended to the data before processing the final chunk. This is a workaround to handle the last word properly, as \n alone might be misinterpreted if the last character is \r.  
Note: This \r\n addition at the end is considered a hack and might need further improvement based on specific use cases.
