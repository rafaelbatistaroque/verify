package verify

import (
	"reflect"
	"testing"
)

// Verifier encapsula o valor e o estado do teste
type Verifier struct {
	t       *testing.T
	value   interface{}
	message string
}

// Should inicializa o Verifier
func Should(t *testing.T, value interface{}) *Verifier {
	return &Verifier{t: t, value: value}
}

// Message define uma mensagem personalizada
func (v *Verifier) Message(msg string) *Verifier {
	v.message = msg
	return v
}

// Be verifica igualdade estrita
func (v *Verifier) Be(expected interface{}) *Verifier {
	if !reflect.DeepEqual(v.value, expected) {
		message := "Expected values to have the same properties and values"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected %v (type %T), but got %v (type %T)", message, expected, expected, v.value, v.value)
	}
	return v
}

// NotEqual verifica desigualdade
func (v *Verifier) NotEqual(unexpected interface{}) *Verifier {
	if reflect.DeepEqual(v.value, unexpected) {
		message := "Expected values to be different"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected not %v (type %T), but got %v (type %T)", message, unexpected, unexpected, v.value, v.value)
	}
	return v
}

// NotEmpty verifica se uma string não está vazia
func (v *Verifier) NotEmpty() *Verifier {
	str, ok := v.asString()
	if !ok {
		return v
	}

	if str == "" {
		message := "Expected string to be not empty"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: string is empty", message)
	}
	return v
}

// Empty verifica se uma string está vazia
func (v *Verifier) Empty() *Verifier {
	str, ok := v.asString()
	if !ok {
		return v
	}

	if str != "" {
		message := "Expected string to be empty"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: string is not empty", message)
	}
	return v
}

// Nil verifica se o valor é nil
func (v *Verifier) Nil() *Verifier {
	if !isNil(v.value) {
		message := "Expected value to be nil"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected nil, but got %v (type %T)", message, v.value, v.value)
	}
	return v
}

// NotNil verifica se o valor não é nil
func (v *Verifier) NotNil() *Verifier {
	if isNil(v.value) {
		message := "Expected value to not be nil"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected not nil, but got nil", message)
	}
	return v
}

// isNil verifica se um valor é nil
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// BeTrue verifica se o valor booleano é true
func (v *Verifier) BeTrue() *Verifier {
	if value, ok := v.value.(bool); !ok || !value {
		message := "Expected value to be true"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected true, but got %v (type %T)", message, v.value, v.value)
	}
	return v
}

// BeFalse verifica se o valor booleano é false
func (v *Verifier) BeFalse() *Verifier {
	if value, ok := v.value.(bool); !ok || value {
		message := "Expected value to be false"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected false, but got %v (type %T)", message, v.value, v.value)
	}
	return v
}

// Len verifica o comprimento de um slice, mapa, canal, string ou array
func (v *Verifier) Len(expected int) *Verifier {
	if v.value == nil {
		v.t.Errorf("Cannot check length of nil value")
		return v
	}

	val := reflect.ValueOf(v.value)
	if val.Kind() != reflect.Slice &&
		val.Kind() != reflect.Map &&
		val.Kind() != reflect.Chan &&
		val.Kind() != reflect.Array &&
		val.Kind() != reflect.String {
		v.t.Errorf("Expected a slice, map, channel, array, or string, but got %T", v.value)
		return v
	}

	// Comparando o comprimento
	actual := val.Len()
	if actual != expected {
		message := "Expected length to be equal"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected length %d, but got %d", message, expected, actual)
	}
	return v
}

// GT verifica se o valor é maior que o esperado
func (v *Verifier) GT(threshold interface{}) *Verifier {
	val, ok := v.value.(float64)
	if !ok {
		if intVal, isInt := v.value.(int); isInt {
			val = float64(intVal)
		} else {
			v.t.Errorf("Expected a numeric value, but got %v (type %T)", v.value, v.value)
			return v
		}
	}

	thr, ok := threshold.(float64)
	if !ok {
		if intThr, isInt := threshold.(int); isInt {
			thr = float64(intThr)
		} else {
			v.t.Errorf("Expected a numeric threshold, but got %v (type %T)", threshold, threshold)
			return v
		}
	}

	if val <= thr {
		message := "Expected value to be greater than threshold"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected > %v, but got %v", message, thr, val)
	}
	return v
}

// LT verifica se o valor é menor que o esperado
func (v *Verifier) LT(threshold interface{}) *Verifier {
	val, ok := v.value.(float64)
	if !ok {
		if intVal, isInt := v.value.(int); isInt {
			val = float64(intVal)
		} else {
			v.t.Errorf("Expected a numeric value, but got %v (type %T)", v.value, v.value)
			return v
		}
	}

	thr, ok := threshold.(float64)
	if !ok {
		if intThr, isInt := threshold.(int); isInt {
			thr = float64(intThr)
		} else {
			v.t.Errorf("Expected a numeric threshold, but got %v (type %T)", threshold, threshold)
			return v
		}
	}

	if val >= thr {
		message := "Expected value to be lower than threshold"
		if v.message != "" {
			message = v.message
		}
		v.t.Errorf("%s: expected < %v, but got %v", message, thr, val)
	}
	return v
}

// ExpectPanic verifica se uma função causa um pânico
func (v *Verifier) Panic(fn func(), msg ...string) *Verifier {
	defer func() {
		if r := recover(); r == nil {
			message := "Expected function to panic"
			if len(msg) > 0 {
				message = msg[0]
			}
			v.t.Errorf("%s: expected panic, but function did not panic", message)
		}
	}()
	fn()

	return v
}

// Expect NotPanic verifica se uma função não causa um pânico
func (v *Verifier) NotPanic(fn func(), msg ...string) *Verifier {
	defer func() {
		if r := recover(); r != nil {
			message := "Expected function to not panic"
			if len(msg) > 0 {
				message = msg[0]
			}
			v.t.Errorf("%s: expected no panic, but function panicked with %v", message, r)
		}
	}()
	fn()

	return v
}

func (v *Verifier) asString() (string, bool) {
	str, ok := v.value.(string)
	if !ok {
		v.t.Errorf("Expected a string, but got value %v (type %T)", v.value, v.value)
	}
	return str, ok
}
