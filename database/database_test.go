package database

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// simple test example
func TestOpenDatabase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	db, err := Open(path)
	require.NoError(t, err)
	require.NotNil(t, db)
}
