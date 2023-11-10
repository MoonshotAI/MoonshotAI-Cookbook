package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNullableType(t *testing.T) {
	t.Run("string", testNullableType[string]("test", false))
	t.Run("int", testNullableType[int](312, true))
	t.Run("float64", testNullableType[float64](3.12, true))
	t.Run("bool", testNullableType[bool](true, true))
}

func testNullableType[T interface {
	string | int | float64 | bool
}](instance T, isNull bool) func(*testing.T) {
	return func(t *testing.T) {
		serializedVal, err := json.Marshal(instance)
		if err != nil {
			t.Errorf("json.Marshal(%T): %s", instance, err)
			return
		}
		var val NullableType[T]
		if val.IsNull() != isNull {
			t.Errorf("NullableType[%T]: empty value must be null=%v, got null=%v", instance, isNull, !isNull)
			return
		}
		if err = json.Unmarshal([]byte("null"), &val); err != nil {
			t.Errorf("json.Unmarshal(nil): %s", err)
			return
		}
		if val != "" {
			t.Errorf("NullableType[%T]: unmarshal from nil should keep value empty, got %s", instance, val)
			return
		}
		val = "test"
		if err = json.Unmarshal([]byte("null"), &val); err != nil {
			t.Errorf("json.Unmarshal(nil): %s", err)
			return
		}
		if val != "test" {
			t.Errorf("NullableType[%T]: unmarshal from nil must be a noop, got %s", instance, val)
			return
		}
		if err = json.Unmarshal(serializedVal, &val); err != nil {
			t.Errorf("json.Unmarshal(%T): %s", instance, err)
			return
		}
		if val.IsNull() || val == "" {
			t.Errorf("NullableType[%[1]T]: expectes %[1]v, got null/empty value", instance)
			return
		}
		if val.Value() != instance {
			t.Errorf("NullableType[%[2]T]: %[1]s != %[2]v", val, instance)
			return
		}
		var reSerializedVal []byte
		reSerializedVal, err = json.Marshal(val)
		if err != nil {
			t.Errorf("json.Marshal(NullableType[%T]): %s", val, err)
			return
		}
		if !bytes.Equal(reSerializedVal, serializedVal) {
			t.Errorf("NullableType[%T]: serialized %s != %s", instance, string(reSerializedVal), string(serializedVal))
			return
		}
		if isNull {
			val = ""
			reSerializedVal, err = json.Marshal(val)
			if err != nil {
				t.Errorf("json.Marshal(NullableType[%T]): %s", val, err)
				return
			}
			if !bytes.Equal(reSerializedVal, []byte("null")) {
				t.Errorf("NullableType[%T]: serialized %s != null", instance, string(reSerializedVal))
				return
			}
		}
		if err = json.Unmarshal([]byte("[]"), &val); err == nil {
			t.Errorf("json.Unmarshal(\"[]\") expects errors, got nothing")
			return
		}
		if err = json.Unmarshal([]byte("{}"), &val); err == nil {
			t.Errorf("json.Unmarshal(\"{}\") expects errors, got nothing")
			return
		}
	}
}

func TestTokenType(t *testing.T) {
	t.Run("null", testTokenType(nil, "null"))
	t.Run("array", testTokenType(json.Delim('['), "array"))
	t.Run("object", testTokenType(json.Delim('{'), "object"))
	t.Run("bool", testTokenType(true, "bool"))
	t.Run("number", testTokenType(json.Number("0313"), "number"))
	t.Run("string", testTokenType("test", "string"))
	t.Run("unknown", testTokenType(struct{}{}, "unknown"))
}

func testTokenType(token json.Token, expect string) func(*testing.T) {
	return func(t *testing.T) {
		if tokType := tokenType(token); tokType != expect {
			t.Errorf("tokenType: %s != %s", tokType, expect)
			return
		}
	}
}
