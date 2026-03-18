package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

func TestZlibBase64Codec_RoundTrip(t *testing.T) {
	c := zlibBase64Codec{}
	in := []*commonpb.Payload{{
		Metadata: map[string][]byte{
			converter.MetadataEncoding: []byte("json/plain"),
			"encryption-key-id":        []byte("key-any-value"),
		},
		Data: []byte(`{"hello":"temporal","n":42}`),
	}}

	encoded, err := c.Encode(in)
	require.NoError(t, err)
	require.Equal(t, codecEncoding, string(encoded[0].Metadata[converter.MetadataEncoding]))
	require.NotEmpty(t, encoded[0].Metadata[metaOriginalEncoding])
	require.NotEqual(t, string(in[0].Data), string(encoded[0].Data))

	decoded, err := c.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, "json/plain", string(decoded[0].Metadata[converter.MetadataEncoding]))
	require.Equal(t, "key-any-value", string(decoded[0].Metadata["encryption-key-id"]))
	require.Equal(t, string(in[0].Data), string(decoded[0].Data))
	require.NotContains(t, decoded[0].Metadata, metaOriginalEncoding)
}

func TestZlibBase64Codec_DecodePassthrough(t *testing.T) {
	c := zlibBase64Codec{}
	in := []*commonpb.Payload{{
		Metadata: map[string][]byte{converter.MetadataEncoding: []byte("binary/plain")},
		Data:     []byte("opaque"),
	}}

	out, err := c.Decode(in)
	require.NoError(t, err)
	require.Equal(t, "binary/plain", string(out[0].Metadata[converter.MetadataEncoding]))
	require.Equal(t, "opaque", string(out[0].Data))
}

func TestZlibBase64Codec_DecodeInvalidBase64(t *testing.T) {
	c := zlibBase64Codec{}
	in := []*commonpb.Payload{{
		Metadata: map[string][]byte{converter.MetadataEncoding: []byte(codecEncoding)},
		Data:     []byte("%%%not-base64%%%"),
	}}

	_, err := c.Decode(in)
	require.Error(t, err)
}
