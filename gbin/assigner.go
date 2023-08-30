package gbin

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/lspaccatrosi16/go-cli-tools/set"
	"github.com/lspaccatrosi16/go-cli-tools/stack"
)

type assigner[T any] struct {
	stack stack.Stack[string]
}

func newAssigner[T any]() *assigner[T] {
	return &assigner[T]{}
}

func (a *assigner[T]) trace() string {
	buf := bytes.NewBufferString("/")
	a.stack.Reverse()
	for {
		val, ok := a.stack.Pop()
		if !ok {
			break
		}
		formatted := fmt.Sprintf("%s/", val)
		buf.WriteString(formatted)
	}
	return buf.String()
}

func (a *assigner[T]) assign(decoded *reflect.Value) (*T, error) {
	refType := reflect.TypeOf(*new(T))
	assigned, err := a.visit(refType, decoded)
	if err != nil {
		return nil, err
	}
	converted := assigned.Interface().(T)
	return &converted, nil
}

func scalars() *set.Set[reflect.Kind] {
	sc := set.NewSet[reflect.Kind]()
	sc.Add(reflect.String, reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64, reflect.Uint, reflect.Uint64, reflect.Uint8)
	return sc
}

func (a *assigner[T]) visit(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	if !a.matches(ref, decoded) {
		return nil, fmt.Errorf("type %s does not match reference type of %s", decoded.Kind(), ref.Kind())
	}
	var visited *reflect.Value
	var visitError error
	if scalars().Contains(ref.Kind()) {
		visited, visitError = a.visit_scalar(ref, decoded)
	} else {
		switch ref.Kind() {
		case reflect.Interface:
			visited, visitError = a.visit_interface(ref, decoded)
		case reflect.Map:
			visited, visitError = a.visit_map(ref, decoded)
		case reflect.Struct:
			visited, visitError = a.visit_struct(ref, decoded)
		case reflect.Pointer:
			visited, visitError = a.visit_ptr(ref, decoded)
		case reflect.Slice:
			visited, visitError = a.visit_slice(ref, decoded)
		default:
			return nil, fmt.Errorf("type: %s is not currently supported for serialization", ref.Kind())
		}
	}
	if visitError != nil {
		return nil, visitError
	} else if visited == nil {
		return nil, fmt.Errorf("visited value is nil")
	}
	return visited, nil
}

func (a *assigner[T]) visit_interface(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	decInner := decoded.Elem()
	a.stack.Push("interface")
	visited, err := a.visit(decInner.Type(), &decInner)
	if err != nil {
		return nil, err
	}
	outer := reflect.New(reflect.TypeOf(visited.Interface())).Elem()
	outer.Set(*visited)
	a.stack.Pop()
	return &outer, nil
}

func (a *assigner[T]) visit_map(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	keyType := ref.Key()
	valType := ref.Elem()
	stackEntry := fmt.Sprintf("map[%s]%s", keyType.Kind(), valType.Kind())
	a.stack.Push(stackEntry)
	iter := decoded.MapRange()
	newMap := reflect.MakeMap(ref)
	for {
		if !iter.Next() {
			break
		}
		k := iter.Key()
		kEntry := fmt.Sprintf("key[%v]", k.Interface())
		a.stack.Push(kEntry)
		kVisited, err := a.visit(keyType, &k)
		if err != nil {
			return nil, err
		}
		a.stack.Pop()
		vEntry := fmt.Sprintf("val[%v]", k.Interface())
		a.stack.Push(vEntry)

		v := iter.Value()
		vVisited, err := a.visit(valType, &v)
		if err != nil {
			return nil, err
		}
		a.stack.Pop()
		newMap.SetMapIndex(*kVisited, *vVisited)
	}
	a.stack.Pop()
	return &newMap, nil
}

func (a *assigner[T]) visit_struct(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	a.stack.Push("struct")
	n := decoded.NumField()
	decType := decoded.Type()
	newStruct := reflect.New(ref).Elem()
	for i := 0; i < n; i++ {
		dFieldT := decType.Field(i)
		dFieldV := decoded.Field(i)
		name := dFieldT.Name
		fEntry := fmt.Sprintf("field[%s]", name)
		a.stack.Push(fEntry)
		rField, found := ref.FieldByName(name)
		if !found {
			return nil, fmt.Errorf("decoded struct has field of name %s but not found in reference type", name)
		}
		neededType := rField.Type
		a.stack.Push("val")
		visited, vErr := a.visit(neededType, &dFieldV)
		if vErr != nil {
			return nil, vErr
		}
		a.stack.Pop()
		a.stack.Pop()
		newStruct.FieldByName(name).Set(*visited)
	}
	a.stack.Pop()
	return &newStruct, nil
}

func (a *assigner[T]) visit_ptr(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	a.stack.Push("ptr")
	refPointedAt := ref.Elem()
	decPointedAt := decoded.Elem()
	visited, err := a.visit(refPointedAt, &decPointedAt)
	if err != nil {
		return nil, err
	}
	var ptr reflect.Value
	if visited.CanAddr() {
		ptr = visited.Addr()
	} else {
		ptr = reflect.New(visited.Type())
		ptr.Elem().Set(*visited)
	}
	a.stack.Pop()
	return &ptr, nil
}

func (a *assigner[T]) visit_slice(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	a.stack.Push(fmt.Sprintf("slice[%s]", decoded.Type().Elem().Kind()))
	n := decoded.Len()
	newSlice := reflect.New(ref).Elem()
	for i := 0; i < n; i++ {
		el := decoded.Index(i)
		a.stack.Push(fmt.Sprintf("el%d", i))
		visited, err := a.visit(ref.Elem(), &el)
		if err != nil {
			return nil, err
		}
		newSlice = reflect.Append(newSlice, *visited)
		a.stack.Pop()
	}
	a.stack.Pop()
	return &newSlice, nil
}

func (a *assigner[T]) visit_scalar(ref reflect.Type, decoded *reflect.Value) (*reflect.Value, error) {
	a.stack.Push(fmt.Sprintf("scalar[%s]", decoded.Type().Kind()))
	if !decoded.CanConvert(ref) {
		return nil, fmt.Errorf("cannot convert type %s to %s", decoded.Kind(), ref.Kind())
	}
	converted := decoded.Convert(ref)
	newVal := reflect.New(ref).Elem()
	newVal.Set(converted)
	a.stack.Pop()
	return &newVal, nil
}

func (a *assigner[T]) matches(x reflect.Type, y *reflect.Value) bool {
	if (x.Kind() == reflect.Int || x.Kind() == reflect.Int64) && ((*y).Kind() == reflect.Int || (*y).Kind() == reflect.Int64) {
		return true
	}
	return x.Kind() == (*y).Kind()
}
