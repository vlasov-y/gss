package config

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
	"gopkg.in/yaml.v3"
)

func DecodeRoot(f reflect.Type, t reflect.Type, data any) (any, error) {
	// Ensure the target type is Headers
	if t != reflect.TypeFor[Root]() {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("expected string but got %T", data)
	}
	root := data.(string)
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve absolute path: %w", err)
	}
	stat, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("failed to stat: %w", err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("root is not a directory: %s", root)
	}
	return Root(filepath.Clean(root)), nil
}

// DecodeHeaders decodes data into http.Header type.
// Supports string (YAML), map[string]any, map[string]string, and map[string][]string as input.
func DecodeHeaders(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	// Ensure the target type is http.Header
	if t != reflect.TypeOf(http.Header{}) {
		return data, nil
	}
	switch f {
	case reflect.TypeFor[string]():
		return decodeHeadersFromYAML(data.(string))
	case reflect.TypeFor[map[string]any]():
		return decodeHeadersFromMap(data.(map[string]any))
	case reflect.TypeFor[map[string]string]():
		return decodeHeadersFromStringMap(data.(map[string]string))
	case reflect.TypeFor[map[string][]string]():
		return decodeHeadersFromMultiValueMap(data.(map[string][]string))
	default:
		return nil, fmt.Errorf("unsupported headers type: %T", data)
	}
}

// decodeHeadersFromYAML unmarshals a YAML string into http.Header.
func decodeHeadersFromYAML(yamlStr string) (http.Header, error) {
	var rawHeaders map[string]any
	if err := yaml.Unmarshal([]byte(yamlStr), &rawHeaders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return decodeHeadersFromMap(rawHeaders)
}

// decodeHeadersFromMap handles map[string]any and converts it into http.Header.
func decodeHeadersFromMap(m map[string]any) (http.Header, error) {
	headers := http.Header{}
	for key, value := range m {
		switch v := value.(type) {
		case string:
			if err := addHeader(&headers, key, value.(string)); err != nil {
				return nil, err
			}
		case []string:
			for _, item := range value.([]string) {
				if err := addHeader(&headers, key, item); err != nil {
					return nil, err
				}
			}
		case []any:
			for i, item := range value.([]any) {
				if str, ok := item.(string); ok {
					if err := addHeader(&headers, key, str); err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("invalid header value for key '%s' at index '%d': could not cast to string: %v", key, i, item)
				}
			}
		default:
			return nil, fmt.Errorf("invalid header value for key '%s': expected string or []string, got %T", key, v)
		}
	}
	return headers, nil
}

// decodeHeadersFromStringMap handles map[string]string and converts it into http.Header.
func decodeHeadersFromStringMap(m map[string]string) (http.Header, error) {
	headers := http.Header{}
	for key, value := range m {
		if err := addHeader(&headers, key, value); err != nil {
			return nil, err
		}
	}
	return headers, nil
}

// decodeHeadersFromMultiValueMap handles map[string][]string and converts it into http.Header.
func decodeHeadersFromMultiValueMap(m map[string][]string) (http.Header, error) {
	headers := http.Header{}
	for key, values := range m {
		for _, value := range values {
			if err := addHeader(&headers, key, value); err != nil {
				return nil, err
			}
		}
	}
	return headers, nil
}

// addHeader checks header key for validity and ads it to http.Header
func addHeader(headers *http.Header, key string, value string) error {
	if !httpguts.ValidHeaderFieldName(key) {
		return fmt.Errorf("invalid header key: %s", key)
	}
	headers.Add(key, value)
	return nil
}

func DecodeCompression(_ reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(Compression(0)) {
		return data, nil
	}
	switch v := data.(type) {
	case string:
		// Trim and convert to lowercase for case-insensitive comparisons
		level := strings.ToLower(strings.TrimSpace(v))
		// Map of valid string-based compression levels
		levels := map[string]int8{
			"none":    gzip.NoCompression,
			"default": gzip.DefaultCompression,
			"speed":   gzip.BestSpeed,
			"best":    gzip.BestCompression,
		}
		// Check if level is in the map
		if num, found := levels[level]; found {
			return Compression(num), nil
		}
		// If not found in map, check if it's a valid integer within the gzip range
		if num, err := strconv.Atoi(level); err == nil {
			if num < gzip.DefaultCompression || num > gzip.BestCompression {
				return nil, fmt.Errorf("unsupported compression level: %d (valid range: %d to %d)", num, gzip.DefaultCompression, gzip.BestCompression)
			}
			return Compression(num), nil
		}
		// Return error if the level is invalid
		return nil, fmt.Errorf("unsupported compression level: %s", level)

	case int, int8, int16, int32, int64:
		// Handle integer types directly
		num := reflect.ValueOf(data).Int()
		if num < int64(gzip.DefaultCompression) || num > int64(gzip.BestCompression) {
			return nil, fmt.Errorf("unsupported compression level: %d (valid range: %d to %d)", num, gzip.DefaultCompression, gzip.BestCompression)
		}
		return Compression(num), nil

	case uint, uint8, uint16, uint32, uint64:
		// Handle unsigned integers as well
		num := reflect.ValueOf(data).Uint()
		if num > uint64(gzip.BestCompression) {
			return nil, fmt.Errorf("unsupported compression level: %d (valid range: %d to %d)", num, gzip.NoCompression, gzip.BestCompression)
		}
		return Compression(num), nil

	default:
		// Unsupported type
		return nil, fmt.Errorf("unsupported compression level type: %T", data)
	}
}
