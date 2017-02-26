package server

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

const cacheVersion = 1 // Increment this to bust ETag cache

// QueryETagMatcher is a middleware which handles conditional GETs by hashing
// the query and comparing it with the ETag.
func QueryETagMatcher(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		etag := buildETagFromRequest(r)
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func buildETagFromRequest(r *http.Request) string {
	return md5sum(r.URL.Query().Encode())
}

func md5sum(query string) string {
	hasher := sha256.New()
	hasher.Write([]byte(query))
	hasher.Write([]byte{cacheVersion})
	return hex.EncodeToString(hasher.Sum(nil))
}
