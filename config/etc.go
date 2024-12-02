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

func DecodeRoot(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
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
// Supports YAML strings or map[string]interface{} as input.
func DecodeHeaders(_ reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	// Ensure the target type is http.Header
	if t != reflect.TypeOf(http.Header{}) {
		return data, nil
	}
	headers := http.Header{}
	switch v := data.(type) {
	// Case 1: Data is a string (e.g., from an environment variable)
	case string:
		yamlStr := v

		// Unmarshal YAML to map[string]interface{}
		var rawHeaders map[string]interface{}
		if err := yaml.Unmarshal([]byte(yamlStr), &rawHeaders); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
		}

		// Validate and add to headers
		for key, value := range rawHeaders {
			if !httpguts.ValidHeaderFieldName(key) {
				return nil, fmt.Errorf("invalid header key: %s", key)
			}
			switch val := value.(type) {
			case string:
				headers.Add(key, val)
			case []interface{}:
				for _, item := range val {
					strItem, ok := item.(string)
					if !ok {
						return nil, fmt.Errorf("invalid value in list for key '%s': expected string, got %T", key, item)
					}
					headers.Add(key, strItem)
				}
			default:
				return nil, fmt.Errorf("invalid header value for key '%s': expected string or []string, got %T", key, value)
			}
		}

	// Case 2: Data is already a map[string]interface{} (e.g., from a YAML file)
	case map[string]interface{}:
		for key, value := range v {
			if !httpguts.ValidHeaderFieldName(key) {
				return nil, fmt.Errorf("invalid header key: %s", key)
			}
			switch val := value.(type) {
			case string:
				headers.Add(key, val)
			case []interface{}:
				for _, item := range val {
					strItem, ok := item.(string)
					if !ok {
						return nil, fmt.Errorf("invalid value in list for key '%s': expected string, got %T", key, item)
					}
					headers.Add(key, strItem)
				}
			default:
				return nil, fmt.Errorf("invalid header value for key '%s': expected string or []string, got %T", key, value)
			}
		}

	// Case 3: Unsupported type
	default:
		return nil, fmt.Errorf("unsupported headers type: %T", data)
	}

	return headers, nil
}

func DecodeCompression(_ reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
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
