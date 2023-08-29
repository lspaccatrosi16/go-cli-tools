package gbin

import (
	"fmt"
	"reflect"
)

type assigner[T any] struct {
}

func newAssigner[T any]() *assigner[T] {
	return &assigner[T]{}
}

type PreMap = map[reflect.Kind]int

func (a *assigner[T]) assign(decoded *reflect.Value) (*T, error) {
	refType := reflect.TypeOf(*new(T))
	assigned, err := a.visit(refType, decoded)
	if err != nil {
		return nil, err
	}
	converted := assigned.Interface().(T)
	return &converted, nil
}

func (a *assigner[T]) visit(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	if decoded.Kind() == reflect.Interface {
		val := decoded.Elem()
		return a.visit(ref, &val)
	}
	if !a.matches(ref, decoded) {
		return nil, fmt.Errorf("type %s does not match reference type of %s", decoded.Kind(), ref.Kind())
	}
	var visited *reflect.Value
	var visitError error
	switch ref.Kind() {
	case reflect.Map:
		visited, visitError = a.visit_map(ref, decoded)
	case reflect.Struct:
		visited, visitError = a.visit_struct(ref, decoded)
	case reflect.Pointer:
		visited, visitError = a.visit_ptr(ref, decoded)
	case reflect.Slice:
		visited, visitError = a.visit_slice(ref, decoded)
	case reflect.String:
		visited, visitError = a.visit_scalar(ref, decoded)
	case reflect.Bool:
		visited, visitError = a.visit_scalar(ref, decoded)
	case reflect.Int, reflect.Int64:
		visited, visitError = a.visit_scalar(ref, decoded)
	case reflect.Float64:
		visited, visitError = a.visit_scalar(ref, decoded)
	default:
		return nil, fmt.Errorf("type: %s is not currently supported for serialization", ref.Kind())
	}
	if visitError != nil {
		return nil, visitError
	} else if visited == nil {
		return nil, fmt.Errorf("visited value is nil")
	}
	return visited, nil
}

func (a *assigner[T]) visit_map(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	keyType := ref.Key()
	valType := ref.Elem()
	iter := decoded.MapRange()
	newMap := reflect.MakeMap(ref)
	for {
		if !iter.Next() {
			break
		}
		k := iter.Key()
		kVisited, err := a.visit(keyType, &k)
		if err != nil {
			return nil, err
		}
		v := iter.Value()
		vVisited, err := a.visit(valType, &v)
		if err != nil {
			return nil, err
		}
		newMap.SetMapIndex(*kVisited, *vVisited)
	}
	return &newMap, nil
}

func (a *assigner[T]) visit_struct(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	n := decoded.NumField()
	decType := decoded.Type()
	newStruct := reflect.New(ref).Elem()
	for i := 0; i < n; i++ {
		dFieldT := decType.Field(i)
		dFieldV := decoded.Field(i)
		name := dFieldT.Name
		rField, found := ref.FieldByName(name)
		if !found {
			return nil, fmt.Errorf("decoded struct has field of name %s but not found in reference type", name)
		}
		neededType := rField.Type
		visited, vErr := a.visit(neededType, &dFieldV)
		if vErr != nil {
			return nil, vErr
		}
		newStruct.FieldByName(name).Set(*visited)
	}
	return &newStruct, nil
}

func (a *assigner[T]) visit_ptr(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	refPointedAt := ref.Elem()
	decPointedAt := decoded.Elem()
	visited, err := a.visit(refPointedAt, &decPointedAt)
	if err != nil {
		return nil, err
	}
	ptr := visited.Addr()
	return &ptr, nil
}

func (a *assigner[T]) visit_slice(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	n := decoded.Len()
	newSlice := reflect.New(ref).Elem()
	for i := 0; i < n; i++ {
		el := decoded.Index(i)
		visited, err := a.visit(ref.Elem(), &el)
		if err != nil {
			return nil, err
		}
		newSlice = reflect.Append(newSlice, *visited)
	}
	return &newSlice, nil
}

func (a *assigner[T]) visit_scalar(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	if !decoded.CanConvert(ref) {
		return nil, fmt.Errorf("cannot convert type %s to %s", decoded.Kind(), ref.Kind())
	}
	converted := decoded.Convert(ref)
	newVal := reflect.New(ref).Elem()
	newVal.Set(converted)
	return &newVal, nil
}

func (a *assigner[T]) matches(x reflect.Type, y *reflect.Value) bool {
	if (x.Kind() == reflect.Int || x.Kind() == reflect.Int64) && ((*y).Kind() == reflect.Int || (*y).Kind() == reflect.Int64) {
		return true
	}
	return x.Kind() == (*y).Kind()
}
