package sequence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Calculator interface defines a mechanism for generating a sequence string.
type Calculator interface {
	Calculate(current interface{}) (interface{}, error)
}

// IncrementCalculator adds 1 to the current sequence (if it's numeric).
// Defaults to 1 if no sequence currently exists.
type IncrementCalculator struct{}

func (c *IncrementCalculator) Calculate(current interface{}) (interface{}, error) {
	if current == nil || current == "" {
		return 1, nil
	}
	switch v := current.(type) {
	case int:
		return v + 1, nil
	case string:
		// Attempt numeric parse, otherwise fail
		var i int
		_, err := fmt.Sscanf(v, "%d", &i)
		if err != nil {
			return nil, fmt.Errorf("cannot increment non-numeric sequence: %v", current)
		}
		return i + 1, nil
	default:
		return nil, fmt.Errorf("unsupported type for increment: %T", current)
	}
}

// FileHashCalculator hashes files or directories and truncates to a specified length.
type FileHashCalculator struct {
	Targets []string
	Length  int
}

func (c *FileHashCalculator) Calculate(current interface{}) (interface{}, error) {
	if len(c.Targets) == 0 {
		return nil, fmt.Errorf("no targets specified for filehash")
	}

	hash := sha256.New()
	var files []string

	// Walk targets to find all files
	for _, target := range c.Targets {
		info, err := os.Stat(target)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %q: %w", target, err)
		}

		if info.IsDir() {
			err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to walk %q: %w", target, err)
			}
		} else {
			files = append(files, target)
		}
	}

	// Ensure determinism
	sort.Strings(files)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %q: %w", file, err)
		}
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to hash file %q: %w", file, err)
		}
		f.Close()
	}

	fullHash := hex.EncodeToString(hash.Sum(nil))
	if c.Length > 0 && c.Length < len(fullHash) {
		return fullHash[:c.Length], nil
	}
	return fullHash, nil
}

// PassedValueCalculator simply returns a literal passed string.
type PassedValueCalculator struct {
	Value string
}

func (c *PassedValueCalculator) Calculate(current interface{}) (interface{}, error) {
	return c.Value, nil
}
