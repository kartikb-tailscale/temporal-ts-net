package app

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

type CodecServer struct {
	Endpoint string

	server   *http.Server
	listener net.Listener
	stopOnce sync.Once
}

const (
	codecEncoding        = "binary/zlib-base64"
	metaOriginalEncoding = "x-original-encoding"
)

type zlibBase64Codec struct{}

func (zlibBase64Codec) Encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		out := clonePayload(p)
		if out.Metadata == nil {
			out.Metadata = make(map[string][]byte, 2)
		}

		origEncoding := append([]byte(nil), out.Metadata[converter.MetadataEncoding]...)
		out.Metadata[metaOriginalEncoding] = origEncoding

		compressed, err := zlibCompress(out.Data)
		if err != nil {
			return payloads, err
		}
		out.Data = []byte(base64.StdEncoding.EncodeToString(compressed))
		out.Metadata[converter.MetadataEncoding] = []byte(codecEncoding)
		result[i] = out
	}
	return result, nil
}

func (zlibBase64Codec) Decode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		out := clonePayload(p)
		if string(out.Metadata[converter.MetadataEncoding]) != codecEncoding {
			result[i] = out
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(string(out.Data))
		if err != nil {
			return payloads, err
		}
		plain, err := zlibDecompress(decoded)
		if err != nil {
			return payloads, err
		}
		out.Data = plain

		if orig, ok := out.Metadata[metaOriginalEncoding]; ok {
			if len(orig) == 0 {
				delete(out.Metadata, converter.MetadataEncoding)
			} else {
				out.Metadata[converter.MetadataEncoding] = append([]byte(nil), orig...)
			}
			delete(out.Metadata, metaOriginalEncoding)
		}
		result[i] = out
	}
	return result, nil
}

func StartCodecServer(port int) (*CodecServer, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, fmt.Errorf("listen codec server: %w", err)
	}

	h := converter.NewPayloadCodecHTTPHandler(zlibBase64Codec{})
	h = withCORS(h)

	srv := &http.Server{Handler: h}
	codec := &CodecServer{
		Endpoint: fmt.Sprintf("http://%s", ln.Addr().String()),
		server:   srv,
		listener: ln,
	}

	go func() {
		err := srv.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			_ = ln.Close()
		}
	}()

	return codec, nil
}

func clonePayload(p *commonpb.Payload) *commonpb.Payload {
	if p == nil {
		return nil
	}
	out := &commonpb.Payload{}
	if len(p.Data) > 0 {
		out.Data = append([]byte(nil), p.Data...)
	}
	if len(p.Metadata) > 0 {
		out.Metadata = make(map[string][]byte, len(p.Metadata))
		for k, v := range p.Metadata {
			out.Metadata[k] = append([]byte(nil), v...)
		}
	}
	return out
}

func zlibCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func zlibDecompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (s *CodecServer) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.server.Shutdown(ctx)
		_ = s.listener.Close()
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Namespace")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
