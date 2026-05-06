// Package baseline provides functionality to save and load a comparison
// result as a baseline, enabling drift detection between runs.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envcmp/internal/comparator"
)

// Record wraps a comparator.Result with metadata for baseline tracking.
type Record struct {
	CreatedAt time.Time          `json:"created_at"`
	LeftFile  string             `json:"left_file"`
	RightFile string             `json:"right_file"`
	Result    comparator.Result  `json:"result"`
}

// Save writes the given Result as a baseline record to the specified file path.
// It returns an error if the file cannot be created or the data cannot be encoded.
func Save(path, leftFile, rightFile string, result comparator.Result) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", path, err)
	}
	defer f.Close()

	rec := Record{
		CreatedAt: time.Now().UTC(),
		LeftFile:  leftFile,
		RightFile: rightFile,
		Result:    result,
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rec); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// Load reads a baseline record from the specified file path.
// It returns an error if the file cannot be opened or the data cannot be decoded.
func Load(path string) (Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return Record{}, fmt.Errorf("baseline: open %q: %w", path, err)
	}
	defer f.Close()

	var rec Record
	if err := json.NewDecoder(f).Decode(&rec); err != nil {
		return Record{}, fmt.Errorf("baseline: decode: %w", err)
	}
	return rec, nil
}
