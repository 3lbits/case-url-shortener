package main

import (
	"context"
	"io"
	"os"
	"testing"
)

// Note to candidates: This is a sample test that should pass, but you can change anything in this file.
func TestLinkFile_Long(t *testing.T) {
	links := &LinkFile{
		Filename: tempLinkFile(t, "e https://elbits.no/\n"),
	}

	got, err := links.Long(context.Background(), "e")
	if err != nil {
		t.Errorf("LinkFile.Long() error = %v; want no error", err)
	}
	if want := "https://elbits.no/"; err == nil && got.String() != want {
		t.Errorf("LinkFile.Long() = %v, want %v", got.String(), want)
	}
}

// Note to candidates: You can use this function to create a temporary file for testing.
func tempLinkFile(t testing.TB, contents string) string {
	t.Helper()
	n, err := os.CreateTemp("", "linkfile*")
	if err != nil {
		t.Fatal(err)
	}
	name := n.Name()
	t.Cleanup(func() {
		if err := os.Remove(n.Name()); err != nil {
			t.Fatalf("failed to remove test file %s", name)
		}
	})

	if _, err := io.WriteString(n, contents); err != nil {
		t.Fatalf("failed to write to test file %s", name)
	}

	return name
}
