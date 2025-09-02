package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	str := String("key", "val")
	assert.Equal(t, StringType, str.typ)
	assert.Equal(t, "key", str.key)
	assert.Equal(t, "val", str.str)

	num := Num("num", 11)
	assert.Equal(t, NumType, num.typ)
	assert.Equal(t, "num", num.key)
	assert.Equal(t, 11, num.num)

	fl := Float("float", 1.1)
	assert.Equal(t, FloatType, fl.typ)
	assert.Equal(t, "float", fl.key)
	assert.Equal(t, 1.1, fl.flt)

	b := Bool("boolean", true)
	assert.Equal(t, BoolType, b.typ)
	assert.Equal(t, "boolean", b.key)
	assert.Equal(t, true, b.b)

	m := make(map[string]any)
	m["object"] = "value"
	an := Any("anything", m)
	assert.Equal(t, AnyType, an.typ)
	assert.Equal(t, "anything", an.key)
	assert.Equal(t, m, an.any)

	err := errors.New("oops")
	er := Error(err)
	assert.Equal(t, ErrorType, er.typ)
	assert.Equal(t, "error", er.key)
	assert.ErrorIs(t, er.err, err)
}

// Test for Query function
func TestQuery(t *testing.T) {
	sql := "select * from t where foo = ? and bar = ?"
	values := []interface{}{0, 1}
	expected := Log{
		typ: AnyType,
		key: "query",
		any: map[string]interface{}{
			"statement": sql,
			"values":    []interface{}{0, 1},
		},
	}

	result := Query(sql, values...)

	assert.Equal(t, expected, result)
}

// Test for Request function
func TestRequest(t *testing.T) {
	value := map[string]interface{}{"hello": "world"}
	expected := Log{
		typ: AnyType,
		key: "request",
		any: value,
	}

	result := Request(value)

	assert.Equal(t, expected, result)
}

// Test for Response function
func TestResponse(t *testing.T) {
	value := map[string]interface{}{"hello": "world"}
	expected := Log{
		typ: AnyType,
		key: "response",
		any: value,
	}

	result := Response(value)

	assert.Equal(t, expected, result)
}

// Test for Header function
func TestHeader(t *testing.T) {
	value := map[string]interface{}{"content-type": "application/json"}
	expected := Log{
		typ: AnyType,
		key: "header",
		any: value,
	}

	result := Header(value)

	assert.Equal(t, expected, result)
}
