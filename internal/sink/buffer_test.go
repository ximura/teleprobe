package sink_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ximura/teleprobe/internal/sink"
)

func TestBuffer_AppendAndFlush(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "buffer-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // clean up

	const maxBytes = 64
	buf, err := sink.NewBuffer(tmpFile.Name(), maxBytes)
	require.NoError(t, err)

	// Append two lines (should not auto-flush)
	err = buf.Append("line 1\n")
	require.NoError(t, err)
	err = buf.Append("line 2\n")
	require.NoError(t, err)

	// Nothing should be flushed yet
	stat, _ := tmpFile.Stat()
	require.Equal(t, int64(0), stat.Size())

	// Force flush
	err = buf.Flush()
	require.NoError(t, err)

	// Verify flushed content
	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	// file should be updated
	stat, _ = tmpFile.Stat()
	require.Equal(t, int64(14), stat.Size())

	require.Equal(t, string(data), "line 1\nline 2\n")
}

func TestBuffer_AutoFlushWhenFull(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "buffer-full-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	const maxBytes = 16 // very small buffer to trigger flush
	buf, err := sink.NewBuffer(tmpFile.Name(), maxBytes)
	require.NoError(t, err)

	// This line should fit, second line should force flush
	err = buf.Append("1234567890\n")
	require.NoError(t, err)

	err = buf.Append("abcdefghij\n") // triggers flush
	require.NoError(t, err)

	// first line should be flushed to file
	stat, _ := tmpFile.Stat()
	require.Equal(t, int64(11), stat.Size())

	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	require.Equal(t, string(data), "1234567890\n")

	// Force flush
	err = buf.Flush()
	require.NoError(t, err)

	// file size should increase after second flush
	stat, _ = tmpFile.Stat()
	require.Equal(t, int64(22), stat.Size())

	data, err = os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, string(data), "1234567890\nabcdefghij\n")
}
