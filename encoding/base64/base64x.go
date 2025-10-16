package base64

import (
	"github.com/cloudwego/base64x"
)

// Encoder provides high-performance base64 encoding capabilities using the cloudwego/base64x library.
//
// Arguments:
// - None
//
// Returns:
// - A configured Encoder instance ready for encoding operations.
type Encoder struct {
	encoding *base64x.Encoding
}

// Decoder provides high-performance base64 decoding capabilities using the cloudwego/base64x library.
//
// Arguments:
// - None
//
// Returns:
// - A configured Decoder instance ready for decoding operations.
type Decoder struct {
	encoding *base64x.Encoding
}

// NewEncoder creates a new high-performance base64 encoder using the standard base64 alphabet.
//
// Arguments:
// - None
//
// Returns:
// - *Encoder: A new encoder instance configured with standard base64 encoding.
func NewEncoder() *Encoder {
	return &Encoder{
		encoding: base64x.StdEncoding,
	}
}

// NewURLEncoder creates a new high-performance base64 encoder using the URL-safe base64 alphabet.
//
// Arguments:
// - None
//
// Returns:
// - *Encoder: A new encoder instance configured with URL-safe base64 encoding.
func NewURLEncoder() *Encoder {
	return &Encoder{
		encoding: base64x.URLEncoding,
	}
}

// NewDecoder creates a new high-performance base64 decoder using the standard base64 alphabet.
//
// Arguments:
// - None
//
// Returns:
// - *Decoder: A new decoder instance configured with standard base64 decoding.
func NewDecoder() *Decoder {
	return &Decoder{
		encoding: base64x.StdEncoding,
	}
}

// NewURLDecoder creates a new high-performance base64 decoder using the URL-safe base64 alphabet.
//
// Arguments:
// - None
//
// Returns:
// - *Decoder: A new decoder instance configured with URL-safe base64 decoding.
func NewURLDecoder() *Decoder {
	return &Decoder{
		encoding: base64x.URLEncoding,
	}
}

// EncodeToString encodes the input byte slice to a base64 encoded string.
//
// Arguments:
// - src: The byte slice to encode.
//
// Returns:
// - string: The base64 encoded string representation of the input.
func (e *Encoder) EncodeToString(src []byte) string {
	return e.encoding.EncodeToString(src)
}

// Encode encodes the source byte slice and writes the result to the destination slice.
//
// Arguments:
// - dst: The destination byte slice where encoded data will be written.
// - src: The source byte slice to encode.
//
// Returns:
// - None
func (e *Encoder) Encode(dst, src []byte) {
	e.encoding.Encode(dst, src)
}

// EncodedLen calculates the length of the base64 encoding of n source bytes.
//
// Arguments:
// - n: The number of source bytes.
//
// Returns:
// - int: The length of the encoded base64 output.
func (e *Encoder) EncodedLen(n int) int {
	return e.encoding.EncodedLen(n)
}

// DecodeString decodes a base64 encoded string to its original byte representation.
//
// Arguments:
// - s: The base64 encoded string to decode.
//
// Returns:
// - []byte: The decoded byte slice.
// - error: An error if the input string is not valid base64.
func (d *Decoder) DecodeString(s string) ([]byte, error) {
	return d.encoding.DecodeString(s)
}

// Decode decodes the source byte slice and writes the result to the destination slice.
//
// Arguments:
// - dst: The destination byte slice where decoded data will be written.
// - src: The source byte slice containing base64 encoded data.
//
// Returns:
// - int: The number of bytes written to dst.
// - error: An error if the input is not valid base64.
func (d *Decoder) Decode(dst, src []byte) (int, error) {
	return d.encoding.Decode(dst, src)
}

// DecodedLen calculates the maximum length of the decoded data corresponding to n base64 bytes.
//
// Arguments:
// - n: The number of base64 encoded bytes.
//
// Returns:
// - int: The maximum length of the decoded output.
func (d *Decoder) DecodedLen(n int) int {
	return d.encoding.DecodedLen(n)
}

// EncodeBytes is a convenience function that encodes a byte slice to base64 using standard encoding.
//
// Arguments:
// - src: The byte slice to encode.
//
// Returns:
// - []byte: The base64 encoded byte slice.
func EncodeBytes(src []byte) []byte {
	dst := make([]byte, base64x.StdEncoding.EncodedLen(len(src)))
	base64x.StdEncoding.Encode(dst, src)
	return dst
}

// DecodeBytes is a convenience function that decodes a base64 encoded byte slice using standard encoding.
//
// Arguments:
// - src: The base64 encoded byte slice to decode.
//
// Returns:
// - []byte: The decoded byte slice.
// - error: An error if the input is not valid base64.
func DecodeBytes(src []byte) ([]byte, error) {
	dst := make([]byte, base64x.StdEncoding.DecodedLen(len(src)))
	n, err := base64x.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// EncodeString is a convenience function that encodes a byte slice to a base64 string using standard encoding.
//
// Arguments:
// - src: The byte slice to encode.
//
// Returns:
// - string: The base64 encoded string.
func EncodeString(src []byte) string {
	return base64x.StdEncoding.EncodeToString(src)
}

// DecodeString is a convenience function that decodes a base64 string using standard encoding.
//
// Arguments:
// - s: The base64 encoded string to decode.
//
// Returns:
// - []byte: The decoded byte slice.
// - error: An error if the input string is not valid base64.
func DecodeString(s string) ([]byte, error) {
	return base64x.StdEncoding.DecodeString(s)
}
