package geo

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

type Resolver struct {
	db *geoip2.Reader
}

// NewResolver opens a MaxMind GeoLite2-Country database.
// If dbPath is empty, it returns nil, nil (feature disabled).
// If the database file does not exist and downloadURL is set, it downloads the file first.
func NewResolver(dbPath, downloadURL string) (*Resolver, error) {
	if dbPath == "" {
		return nil, nil
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if downloadURL == "" {
			return nil, fmt.Errorf("GeoIP database not found at %s and no DownloadURL configured", dbPath)
		}
		if err := downloadFile(dbPath, downloadURL); err != nil {
			return nil, fmt.Errorf("failed to download GeoIP database: %w", err)
		}
		slog.Info("GeoIP database downloaded", "path", dbPath, "url", downloadURL)
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Resolver{db: db}, nil
}

// downloadFile downloads a file from url and saves it to dest, creating parent directories as needed.
func downloadFile(dest, url string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		os.Remove(dest)
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// Country returns the ISO 3166-1 alpha-2 country code for the given IP.
// It is nil-safe: calling Country on a nil Resolver returns "".
func (r *Resolver) Country(ipStr string) string {
	if r == nil {
		return ""
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	record, err := r.db.Country(ip)
	if err != nil {
		return ""
	}
	return record.Country.IsoCode
}

// Close releases the database resources.
func (r *Resolver) Close() error {
	if r == nil {
		return nil
	}
	return r.db.Close()
}
