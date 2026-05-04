package base64

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncoder_EncodeToString validates that the encoder correctly encodes various inputs to base64 strings.
// This test ensures proper handling of empty strings, standard text, and binary data.
func TestEncoder_EncodeToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		encoder  *Encoder
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: "",
			encoder:  NewEncoder(),
		},
		{
			name:     "Hello World",
			input:    []byte("Hello World"),
			expected: "SGVsbG8gV29ybGQ=",
			encoder:  NewEncoder(),
		},
		{
			name:     "Binary data",
			input:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			expected: "AAECAv/+/Q==",
			encoder:  NewEncoder(),
		},
		{
			name:     "URL encoding with special chars",
			input:    []byte("hello?world&foo=bar"),
			expected: "aGVsbG8_d29ybGQmZm9vPWJhcg==",
			encoder:  NewURLEncoder(),
		},
		{
			name:     "Long text",
			input:    []byte("The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs."),
			expected: "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZy4gUGFjayBteSBib3ggd2l0aCBmaXZlIGRvemVuIGxpcXVvciBqdWdzLg==",
			encoder:  NewEncoder(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.encoder.EncodeToString(tt.input)
			assert.Equal(t, tt.expected, result, "Encoded string should match expected value.")
		})
	}
}

// TestDecoder_DecodeString validates that the decoder correctly decodes base64 strings to their original form.
// This test covers edge cases including empty strings, padding variations, and invalid inputs.
func TestDecoder_DecodeString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []byte
		decoder     *Decoder
		expectError bool
	}{
		{
			name:        "Empty input",
			input:       "",
			expected:    []byte{},
			decoder:     NewDecoder(),
			expectError: false,
		},
		{
			name:        "Hello World",
			input:       "SGVsbG8gV29ybGQ=",
			expected:    []byte("Hello World"),
			decoder:     NewDecoder(),
			expectError: false,
		},
		{
			name:        "Binary data",
			input:       "AAECAv/+/Q==",
			expected:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			decoder:     NewDecoder(),
			expectError: false,
		},
		{
			name:        "URL decoding",
			input:       "aGVsbG8_d29ybGQmZm9vPWJhcg==",
			expected:    []byte("hello?world&foo=bar"),
			decoder:     NewURLDecoder(),
			expectError: false,
		},
		{
			name:        "Invalid base64",
			input:       "This is not base64!",
			expected:    nil,
			decoder:     NewDecoder(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.decoder.DecodeString(tt.input)

			if tt.expectError {
				assert.Error(t, err, "Should return an error for invalid input.")
			} else {
				require.NoError(t, err, "Should not return an error for valid input.")
				assert.Equal(t, tt.expected, result, "Decoded bytes should match expected value.")
			}
		})
	}
}

// TestEncoder_Encode validates the Encode method which writes to a pre-allocated destination buffer.
// This test ensures correct buffer usage and encoding accuracy.
func TestEncoder_Encode(t *testing.T) {
	encoder := NewEncoder()
	src := []byte("Hello, World!")

	// Calculate required destination size
	dstLen := encoder.EncodedLen(len(src))
	dst := make([]byte, dstLen)

	// Encode
	encoder.Encode(dst, src)

	// Verify
	expected := "SGVsbG8sIFdvcmxkIQ=="
	assert.Equal(t, expected, string(dst), "Encoded result should match expected base64 string.")
}

// TestDecoder_Decode validates the Decode method which writes to a pre-allocated destination buffer.
// This test verifies correct buffer usage, return values, and decoding accuracy.
func TestDecoder_Decode(t *testing.T) {
	decoder := NewDecoder()
	src := []byte("SGVsbG8sIFdvcmxkIQ==")

	// Calculate required destination size
	dstLen := decoder.DecodedLen(len(src))
	dst := make([]byte, dstLen)

	// Decode
	n, err := decoder.Decode(dst, src)

	// Verify
	require.NoError(t, err, "Should decode without error.")
	assert.Equal(t, "Hello, World!", string(dst[:n]), "Decoded result should match original text.")
	assert.Equal(t, 13, n, "Should return correct number of decoded bytes.")
}

// TestEncodedLen validates that EncodedLen correctly calculates the output size for various input lengths.
// This ensures proper buffer allocation for encoding operations.
func TestEncodedLen(t *testing.T) {
	encoder := NewEncoder()

	tests := []struct {
		inputLen    int
		expectedLen int
	}{
		{0, 0},
		{1, 4},
		{2, 4},
		{3, 4},
		{4, 8},
		{5, 8},
		{6, 8},
		{100, 136},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			result := encoder.EncodedLen(tt.inputLen)
			assert.Equal(t, tt.expectedLen, result, "EncodedLen should return correct length for input size %d.", tt.inputLen)
		})
	}
}

// TestDecodedLen validates that DecodedLen correctly calculates the maximum output size for various input lengths.
// This ensures proper buffer allocation for decoding operations.
func TestDecodedLen(t *testing.T) {
	decoder := NewDecoder()

	tests := []struct {
		inputLen    int
		expectedLen int
	}{
		{0, 0},
		{4, 3},
		{8, 6},
		{12, 9},
		{136, 102},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			result := decoder.DecodedLen(tt.inputLen)
			assert.Equal(t, tt.expectedLen, result, "DecodedLen should return correct maximum length for input size %d.", tt.inputLen)
		})
	}
}

// TestConvenienceFunctions validates the package-level convenience functions for quick encoding/decoding.
// These functions provide a simple API for common use cases.
func TestConvenienceFunctions(t *testing.T) {
	originalData := []byte("Testing convenience functions!")

	t.Run("EncodeString and DecodeString", func(t *testing.T) {
		encoded := EncodeString(originalData)
		decoded, err := DecodeString(encoded)

		require.NoError(t, err, "DecodeString should not return an error.")
		assert.Equal(t, originalData, decoded, "Round-trip encoding/decoding should preserve original data.")
	})

	t.Run("EncodeBytes and DecodeBytes", func(t *testing.T) {
		encoded := EncodeBytes(originalData)
		decoded, err := DecodeBytes(encoded)

		require.NoError(t, err, "DecodeBytes should not return an error.")
		assert.Equal(t, originalData, decoded, "Round-trip encoding/decoding should preserve original data.")
	})
}

// TestLargeData validates encoding and decoding of large data sets to ensure performance and correctness at scale.
// This test helps identify any issues with buffer management or memory allocation.
func TestLargeData(t *testing.T) {
	// Create a large data set (1MB)
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	encoder := NewEncoder()
	decoder := NewDecoder()

	// Encode
	encoded := encoder.EncodeToString(largeData)

	// Decode
	decoded, err := decoder.DecodeString(encoded)

	require.NoError(t, err, "Should decode large data without error.")
	assert.True(t, bytes.Equal(largeData, decoded), "Large data should round-trip correctly.")
}

// BenchmarkEncoder_EncodeToString measures the performance of string encoding operations.
// This benchmark helps identify performance characteristics across different data sizes.
func BenchmarkEncoder_EncodeToString(b *testing.B) {
	sizes := []int{10, 100, 1024, 10240, 102400}

	for _, size := range sizes {
		b.Run(b.Name(), func(b *testing.B) {
			data := make([]byte, size)
			for i := range data {
				data[i] = byte(i % 256)
			}

			encoder := NewEncoder()
			b.SetBytes(int64(size))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = encoder.EncodeToString(data)
			}
		})
	}
}

// BenchmarkDecoder_DecodeString measures the performance of string decoding operations.
// This benchmark helps identify performance characteristics across different data sizes.
func BenchmarkDecoder_DecodeString(b *testing.B) {
	sizes := []int{10, 100, 1024, 10240, 102400}

	for _, size := range sizes {
		b.Run(b.Name(), func(b *testing.B) {
			data := make([]byte, size)
			for i := range data {
				data[i] = byte(i % 256)
			}

			encoder := NewEncoder()
			encoded := encoder.EncodeToString(data)
			decoder := NewDecoder()

			b.SetBytes(int64(size))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = decoder.DecodeString(encoded)
			}
		})
	}
}

// BenchmarkEncoder_Encode measures the performance of buffer-based encoding operations.
// This benchmark tests the efficiency of pre-allocated buffer usage.
func BenchmarkEncoder_Encode(b *testing.B) {
	size := 1024
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	encoder := NewEncoder()
	dst := make([]byte, encoder.EncodedLen(len(data)))

	b.SetBytes(int64(size))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder.Encode(dst, data)
	}
}

// BenchmarkDecoder_Decode measures the performance of buffer-based decoding operations.
// This benchmark tests the efficiency of pre-allocated buffer usage.
func BenchmarkDecoder_Decode(b *testing.B) {
	size := 1024
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	encoder := NewEncoder()
	encoded := encoder.EncodeToString(data)
	encodedBytes := []byte(encoded)

	decoder := NewDecoder()
	dst := make([]byte, decoder.DecodedLen(len(encodedBytes)))

	b.SetBytes(int64(size))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = decoder.Decode(dst, encodedBytes)
	}
}

// BenchmarkConvenienceFunctions compares the performance of convenience functions against direct usage.
// This helps users understand the overhead of using simplified APIs.
func BenchmarkConvenienceFunctions(b *testing.B) {
	size := 1024
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.Run("EncodeString", func(b *testing.B) {
		b.SetBytes(int64(size))
		for i := 0; i < b.N; i++ {
			_ = EncodeString(data)
		}
	})

	b.Run("DecodeString", func(b *testing.B) {
		encoded := EncodeString(data)
		b.SetBytes(int64(size))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = DecodeString(encoded)
		}
	})
}
