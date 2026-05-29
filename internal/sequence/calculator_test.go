package sequence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIncrementCalculator(t *testing.T) {
	calc := &IncrementCalculator{}

	res, err := calc.Calculate(nil)
	if err != nil || res != 1 {
		t.Fatalf("expected 1, got %v", res)
	}

	res, err = calc.Calculate(42)
	if err != nil || res != 43 {
		t.Fatalf("expected 43, got %v", res)
	}

	res, err = calc.Calculate("10")
	if err != nil || res != 11 {
		t.Fatalf("expected 11, got %v", res)
	}
}

func TestPassedValueCalculator(t *testing.T) {
	calc := &PassedValueCalculator{Value: "custom-build-xyz"}
	res, err := calc.Calculate(nil)
	if err != nil || res != "custom-build-xyz" {
		t.Fatalf("expected custom-build-xyz, got %v", res)
	}
}

func TestFileHashCalculator(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hash-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	f1 := filepath.Join(tmpDir, "A.txt")
	os.WriteFile(f1, []byte("hello"), 0644)

	calc := &FileHashCalculator{
		Targets: []string{tmpDir},
		Length:  7,
	}

	res, err := calc.Calculate(nil)
	if err != nil {
		t.Fatal(err)
	}
	hashStr, ok := res.(string)
	if !ok || len(hashStr) != 7 {
		t.Fatalf("expected 7-char hash, got %v", res)
	}
}
