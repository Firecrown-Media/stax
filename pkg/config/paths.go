package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetValueByPath retrieves a value from the config using dot notation
// Example: "project.name", "wpengine.install", "ddev.php_version"
func GetValueByPath(cfg *Config, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	// Use reflection to navigate the struct
	value := reflect.ValueOf(cfg).Elem()

	for i, part := range parts {
		if !value.IsValid() {
			return nil, fmt.Errorf("invalid path at segment: %s", strings.Join(parts[:i], "."))
		}

		// Handle struct fields
		if value.Kind() == reflect.Struct {
			// Convert dot notation to struct field name (e.g., "php_version" -> "PHPVersion")
			fieldName := toFieldName(part)
			field := value.FieldByName(fieldName)

			if !field.IsValid() {
				return nil, fmt.Errorf("field not found: %s", part)
			}

			value = field
		} else if value.Kind() == reflect.Slice {
			// Handle array/slice indexing
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}

			if index < 0 || index >= value.Len() {
				return nil, fmt.Errorf("index out of range: %d", index)
			}

			value = value.Index(index)
		} else if value.Kind() == reflect.Map {
			// Handle map keys
			mapKey := reflect.ValueOf(part)
			value = value.MapIndex(mapKey)

			if !value.IsValid() {
				return nil, fmt.Errorf("key not found: %s", part)
			}
		} else {
			return nil, fmt.Errorf("cannot traverse non-struct/slice/map type at: %s", part)
		}
	}

	// Return the interface value
	return value.Interface(), nil
}

// SetValueByPath sets a value in the config using dot notation
func SetValueByPath(cfg *Config, path string, valueStr string) error {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("invalid path: %s", path)
	}

	// Use reflection to navigate to the parent of the target field
	value := reflect.ValueOf(cfg).Elem()

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		if value.Kind() != reflect.Struct {
			return fmt.Errorf("cannot traverse non-struct type at: %s", part)
		}

		fieldName := toFieldName(part)
		field := value.FieldByName(fieldName)

		if !field.IsValid() {
			return fmt.Errorf("field not found: %s", part)
		}

		value = field
	}

	// Now set the final field
	lastPart := parts[len(parts)-1]
	fieldName := toFieldName(lastPart)
	field := value.FieldByName(fieldName)

	if !field.IsValid() {
		return fmt.Errorf("field not found: %s", lastPart)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", lastPart)
	}

	// Convert string value to appropriate type
	return setFieldValue(field, valueStr)
}

// setFieldValue sets a reflect.Value from a string representation
func setFieldValue(field reflect.Value, valueStr string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(valueStr)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", valueStr)
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value: %s", valueStr)
		}
		field.SetUint(uintVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(valueStr)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s (use true/false)", valueStr)
		}
		field.SetBool(boolVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", valueStr)
		}
		field.SetFloat(floatVal)

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// toFieldName converts a snake_case or kebab-case string to PascalCase
// Examples: "php_version" -> "PHPVersion", "install" -> "Install"
func toFieldName(name string) string {
	// Handle special cases first
	specialCases := map[string]string{
		"url":               "URL",
		"api":               "API",
		"ssh":               "SSH",
		"ssh_gateway":       "SSHGateway",
		"php_version":       "PHPVersion",
		"mysql_version":     "MySQLVersion",
		"mysql_type":        "MySQLType",
		"nodejs_version":    "NodeJSVersion",
		"npm":               "NPM",
		"phpcs":             "PHPCS",
		"nfs_mount_enabled": "NFSMountEnabled",
		"fqdns":             "FQDNs",
		"additional_fqdns":  "AdditionalFQDNs",
		"ttl":               "TTL",
		"ddev":              "DDEV",
		"wpengine":          "WPEngine",
		"wordpress":         "WordPress",
	}

	if mapped, ok := specialCases[name]; ok {
		return mapped
	}

	// Convert snake_case/kebab-case to PascalCase
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-'
	})

	for i, part := range parts {
		if part == "" {
			continue
		}

		// Check if this part has a special case mapping
		if mapped, ok := specialCases[part]; ok {
			parts[i] = mapped
		} else {
			// Capitalize first letter
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(parts, "")
}

// toYAMLKey converts a PascalCase field name to snake_case YAML key
// Examples: "PHPVersion" -> "php_version", "Install" -> "install"
func toYAMLKey(name string) string {
	// Handle special cases
	specialCases := map[string]string{
		"URL":             "url",
		"SSHGateway":      "ssh_gateway",
		"PHPVersion":      "php_version",
		"MySQLVersion":    "mysql_version",
		"MySQLType":       "mysql_type",
		"NodeJSVersion":   "nodejs_version",
		"NFSMountEnabled": "nfs_mount_enabled",
		"FQDNs":           "fqdns",
		"TTL":             "ttl",
		"DDEV":            "ddev",
		"NPM":             "npm",
		"PHPCS":           "phpcs",
	}

	if mapped, ok := specialCases[name]; ok {
		return mapped
	}

	// Convert PascalCase to snake_case
	var result []rune
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}

	return strings.ToLower(string(result))
}

// ValidatePath checks if a path is valid for the config structure
func ValidatePath(cfg *Config, path string) error {
	_, err := GetValueByPath(cfg, path)
	return err
}

// GetAllPaths returns all valid paths in the config
func GetAllPaths(cfg *Config) []string {
	paths := []string{}
	collectPaths(reflect.ValueOf(cfg).Elem(), "", &paths)
	return paths
}

// collectPaths recursively collects all paths in a struct
func collectPaths(value reflect.Value, prefix string, paths *[]string) {
	if !value.IsValid() {
		return
	}

	switch value.Kind() {
	case reflect.Struct:
		typ := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := value.Field(i)

			// Get YAML tag or convert field name
			yamlTag := field.Tag.Get("yaml")
			fieldName := field.Name

			if yamlTag != "" && yamlTag != "-" {
				// Remove omitempty and other options
				fieldName = strings.Split(yamlTag, ",")[0]
			} else {
				fieldName = toYAMLKey(fieldName)
			}

			path := fieldName
			if prefix != "" {
				path = prefix + "." + fieldName
			}

			// Add this path if it's a basic type
			if isBasicType(fieldValue.Kind()) {
				*paths = append(*paths, path)
			}

			// Recurse for nested structs
			if fieldValue.Kind() == reflect.Struct {
				collectPaths(fieldValue, path, paths)
			}
		}

	case reflect.Slice, reflect.Array:
		// Don't enumerate array elements, just note that this is an array path
		if prefix != "" {
			*paths = append(*paths, prefix+"[i]")
		}
	}
}

// isBasicType returns true if the kind represents a basic value type
func isBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}
