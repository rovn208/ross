package util

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	tmpdir := t.TempDir()
	dir := filepath.Join(tmpdir, "dir")
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}

	err := CreateDirectory(dir)
	require.NoError(t, err)

	notExistDir := filepath.Join(tmpdir, "dir2")
	err = CreateDirectory(notExistDir)
	require.NoError(t, err)
}

func TestIsFileAlreadyExists(t *testing.T) {
	tmpdir := t.TempDir()
	file := filepath.Join(tmpdir, "file")
	if err := os.WriteFile(file, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	isFileAlreadyExists := IsFileAlreadyExists(file)
	require.True(t, isFileAlreadyExists)

	require.False(t, IsFileAlreadyExists(filepath.Join(tmpdir, "file2")))
}
