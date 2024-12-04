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
	if len(strings.TrimSpace(value)) == 0 {
		return fmt.Errorf("empty value for key: %s", key)
	}
	headers.Add(key, value)
	return nil
}

func DecodeCompression(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(Compression(0)) {
		return data, nil
	}
	var strLevel string
	switch f.Kind() {
	case reflect.String:
		strLevel = strings.ToLower(strings.TrimSpace(data.(string)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle integer types and convert to string
		strLevel = strconv.FormatInt(reflect.ValueOf(data).Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Handle unsigned integer types and convert to string
		strLevel = strconv.FormatUint(reflect.ValueOf(data).Uint(), 10)
	default:
		// Unsupported type
		return nil, fmt.Errorf("unsupported compression level type: %T", data)
	}
	// Map of valid string-based compression levels
	levels := map[string]int8{
		"none":    gzip.NoCompression,
		"default": gzip.DefaultCompression,
		"speed":   gzip.BestSpeed,
		"best":    gzip.BestCompression,
	}
	// Check if level is in the map
	var num int8
	var ok bool
	var err error
	if num, ok = levels[strLevel]; !ok {
		// If not found in map, check if it's a valid integer within the gzip range
		var num64 int64
		if num64, err = strconv.ParseInt(strLevel, 10, 8); err != nil {
			return nil, fmt.Errorf("level %d is out int8 range", num64)
		}
		num = int8(num64)
	}
	if num < gzip.DefaultCompression || num > gzip.BestCompression {
		return nil, fmt.Errorf("level %d is out of range (valid range: %d to %d)", num, gzip.DefaultCompression, gzip.BestCompression)
	}
	return Compression(num), nil
}
