package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type OperationType string

var (
	QUERY    OperationType = "query"
	MUTATION OperationType = "mutation"
)

type request struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

type response struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors GraphqlErrors              `json:"errors"`
}

type Argument interface {
	Value() interface{}
	Required() bool
}

type Arguments map[string]Argument

type required struct {
	value interface{}
}

func (r required) Value() interface{} {
	return r.value
}

func (required) Required() bool {
	return true
}

func Required(v interface{}) required {
	return required{v}
}

type optional struct {
	value interface{}
}

func (o optional) Value() interface{} {
	return o.value
}

func (optional) Required() bool {
	return false
}

func Optional(v interface{}) optional {
	return optional{v}
}

func isScalarType(t string) bool {
	switch t {
	case "String", "Int", "Float", "Boolean":
		return true
	default:
		return false
	}
}

func unwrapNestedType(t reflect.Type) reflect.Type {
	for k := t.Kind(); k == reflect.Pointer || k == reflect.Slice; k = t.Kind() {
		t = t.Elem()
	}
	return t
}

func getGraphqlType(v interface{}, required bool) (string, error) {
	var t reflect.Type
	if _v, ok := v.(reflect.Type); ok {
		t = _v
	} else {
		t = reflect.TypeOf(v)
	}
	suffix := ""
	if required {
		suffix = "!"
	}

	switch k := t.Kind(); k {
	case reflect.String:
		return "String" + suffix, nil
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return "Int" + suffix, nil
	case reflect.Float32, reflect.Float64:
		return "Float" + suffix, nil
	case reflect.Bool:
		return "Boolean" + suffix, nil
	case reflect.Array, reflect.Slice:
		innerType, err := getGraphqlType(t.Elem(), false)
		if err != nil {
			return "", err
		}

		if !isScalarType(innerType) {
			return "", errors.New("slice type must be scalar")
		}
		return fmt.Sprintf("[%s%s]", innerType, suffix), nil
	default:
		return "", fmt.Errorf("unknown type %s", k)
	}
}

func fieldName(f reflect.StructField) string {
	if graphqlTag := f.Tag.Get("graphql"); graphqlTag != "" {
		return graphqlTag
	}

	jsonTag, ok := f.Tag.Lookup("json")
	if !ok {
		return f.Name
	}

	if jsonTag == "" || jsonTag == "-" {
		return f.Name
	}

	name := strings.Split(jsonTag, ",")[0]
	if name == "" {
		return f.Name
	}

	return name
}

func generateFields(v interface{}) string {
	var t reflect.Type
	if _v, ok := v.(reflect.Type); ok {
		t = _v
	} else {
		t = reflect.TypeOf(v)
	}
	t = unwrapNestedType(t)

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct, got %s", t.Kind()))
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		field := fieldName(f)

		// Recursively build field string if there are any children.
		if f.Type.Kind() == reflect.Struct {
			field = fmt.Sprintf("%s { %s }", field, generateFields(f.Type))
		} else if f.Type.Kind() == reflect.Slice &&
			f.Type.Elem().Kind() == reflect.Struct {
			field = fmt.Sprintf("%s { %s }", field, generateFields(f.Type.Elem()))
		}

		fields = append(fields, field)
	}

	return strings.Join(fields, " ")
}

func generateQuery(operation OperationType,
	name string,
	schema interface{},
	args Arguments,
) (string, error) {
	var variables, arguments []string

	for name, arg := range args {
		t, err := getGraphqlType(arg.Value(), arg.Required())
		if err != nil {
			return "", fmt.Errorf("failed to convert field %s: %w", name, err)
		}

		variables = append(variables,
			fmt.Sprintf("$%s: %s", name, t))
		arguments = append(arguments, fmt.Sprintf("%s: $%s", name, name))
	}
	fields := generateFields(schema)

	return fmt.Sprintf("%s %s(%s) { %s(%s) { %s } }",
		operation,
		name,
		strings.Join(variables, ", "),
		name,
		strings.Join(arguments, ", "),
		fields), nil
}

func Execute(client *http.Client,
	url string,
	operation OperationType,
	name string,
	schema interface{},
	args Arguments,
) error {
	query, err := generateQuery(operation, name, schema, args)
	if err != nil {
		return err
	}

	_args := make(map[string]interface{}, len(args))
	for key, arg := range args {
		_args[key] = arg.Value
	}

	variables, err := json.Marshal(_args)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&request{
		Query: query, Variables: string(variables)})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("expected status 200, got %d: %s", res.StatusCode, body)
	}

	var wrapper response
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return err
	}
	if len(wrapper.Errors) != 0 {
		return wrapper.Errors
	}

	payload, ok := wrapper.Data[name]
	if !ok {
		return nil
	}

	return json.Unmarshal(payload, schema)
}
