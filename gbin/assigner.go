package gbin

import (
	"fmt"
	"reflect"
)

type assigner[T any] struct {
}

type PreMap = map[reflect.Kind]int

var floatPrecedence PreMap = PreMap{
	reflect.Float32: 0,
	reflect.Float64: 1,
}

var intPrecedence PreMap = PreMap{
	reflect.Int8:   0,
	reflect.Uint8:  1,
	reflect.Int16:  1,
	reflect.Uint16: 2,
	reflect.Int32:  2,
	reflect.Uint32: 3,
	reflect.Int64:  3,
	reflect.Int:    3,
	reflect.Uint64: 4,
	reflect.Uint:   4,
}

func (a *assigner[T]) assign(decoded *reflect.Value) (*T, error) {
	refVal := reflect.ValueOf(*new(T))
	err := a.visit(&refVal, decoded)
	if err != nil {
		return nil, err
	}
	assigned := refVal.Interface().(T)
	return &assigned, nil
}

func (a *assigner[T]) visit(ref *reflect.Value, decoded *reflect.Value) error {
	// if !a.matches(ref, decoded) {
	// 	return fmt.Errorf("expected type %s but got type %s", ref.Kind(), decoded.Kind())
	// }

	switch ref.Kind() {
	case reflect.Map:
		return a.visit_map(ref, decoded)
	case reflect.Struct:
		return a.visit_struct(ref, decoded)

	case reflect.Pointer:
		return a.visit_ptr(ref, decoded)

	case reflect.Slice:
		return a.visit_slice(ref, decoded)

	case reflect.String:
		return a.visit_string(ref, decoded)
	case reflect.Bool:
		return a.visit_bool(ref, decoded)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.visit_float(ref, decoded)
	default:
		return fmt.Errorf("type: %s is not currently supported for serialization", ref.Kind())
	}
}

func (a *assigner[T]) visit_map(ref *reflect.Value, decoded *reflect.Value) error {

	return nil
}

func (a *assigner[T]) visit_struct(ref *reflect.Value, decoded *reflect.Value) error {

	return nil
}

func (a *assigner[T]) visit_ptr(ref *reflect.Value, decoded *reflect.Value) error {

	return nil
}

func (a *assigner[T]) visit_slice(ref *reflect.Value, decoded *reflect.Value) error {
	return nil

}

func (a *assigner[T]) visit_string(ref *reflect.Value, decoded *reflect.Value) error {
	return nil

}

func (a *assigner[T]) visit_bool(ref *reflect.Value, decoded *reflect.Value) error {
	return nil

}

func (a *assigner[T]) visit_int(ref *reflect.Value, decoded *reflect.Value) error {
	refPrec := intPrecedence[ref.Kind()]
	decPrec, ok := intPrecedence[decoded.Kind()]
	if !ok {
		return fmt.Errorf("type %s is not compatable with reference type %s", decoded.Kind(), ref.Kind())
	}
	if refPrec >= decPrec {
		val := decoded.Convert(ref.Type())
		ref.Set(val)
		return nil
	} else {
		return fmt.Errorf("cannot safely convert parsed type %s to reference type %s", decoded.Kind(), ref.Kind())
	}
}

func (a *assigner[T]) visit_float(ref *reflect.Value, decoded *reflect.Value) error {
	refPrec := floatPrecedence[ref.Kind()]
	decPrec, ok := floatPrecedence[decoded.Kind()]
	if !ok {
		return fmt.Errorf("type %s is not compatable with reference type %s", decoded.Kind(), ref.Kind())
	}
	if refPrec >= decPrec {
		val := decoded.Float()
		ref.SetFloat(val)
		return nil
	} else {
		return fmt.Errorf("cannot safely convert parsed type %s to reference type %s", decoded.Kind(), ref.Kind())
	}
}

func (a *assigner[T]) matches(x *reflect.Value, y *reflect.Value) bool {
	return (*x).Kind() == (*y).Kind()
}
