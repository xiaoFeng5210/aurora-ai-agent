package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return "{}", nil
	}

	values := make([]string, 0, len(sa))
	for _, item := range sa {
		values = append(values, quotePostgresArrayValue(item))
	}
	return "{" + strings.Join(values, ",") + "}", nil
}

func (sa *StringArray) Scan(value any) error {
	if sa == nil {
		return fmt.Errorf("StringArray: Scan on nil pointer")
	}

	switch v := value.(type) {
	case nil:
		*sa = StringArray{}
		return nil
	case []byte:
		return sa.scanString(string(v))
	case string:
		return sa.scanString(v)
	default:
		return fmt.Errorf("StringArray: unsupported Scan type %T", value)
	}
}

func (sa *StringArray) scanString(value string) error {
	if value == "" || value == "{}" {
		*sa = StringArray{}
		return nil
	}

	items, err := parsePostgresTextArray(value)
	if err != nil {
		return err
	}

	*sa = StringArray(items)
	return nil
}

func quotePostgresArrayValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return `"` + value + `"`
}

func parsePostgresTextArray(value string) ([]string, error) {
	if len(value) < 2 || value[0] != '{' || value[len(value)-1] != '}' {
		return nil, fmt.Errorf("StringArray: invalid PostgreSQL array literal %q", value)
	}

	body := value[1 : len(value)-1]
	if body == "" {
		return []string{}, nil
	}

	var (
		result  []string
		current strings.Builder
		quoted  bool
		escaped bool
	)

	for _, r := range body {

		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case quoted && r == '\\':
			escaped = true
		case r == '"':
			quoted = !quoted
		case !quoted && r == ',':
			result = append(result, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}

	if escaped || quoted {
		return nil, fmt.Errorf("StringArray: invalid PostgreSQL array literal %q", value)
	}

	result = append(result, current.String())
	return result, nil
}
