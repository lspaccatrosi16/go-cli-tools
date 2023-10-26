package structconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type ConfigOutput[T any] struct {
	exec *func() error
	tree *node
}

func (c *ConfigOutput[T]) Analyze(s *T) {
	if s == nil {
		panic("input is nil pointer")
	}

	tree := makeTree(reflect.ValueOf(s))
	c.tree = tree
}

func (c *ConfigOutput[T]) Traverse() {
	if c.tree == nil {
		panic("must analyze before traverse")
	}
	exec := traverseTree(c.tree)
	c.exec = &exec
}

func (c *ConfigOutput[T]) Execute() error {
	if c.exec == nil {
		panic("must traverse before execute")
	}
	return (*c.exec)()
}

func (c *ConfigOutput[T]) Debug() {
	if c.tree != nil {
		fmt.Println(c.tree.String())
	}
}

func (c *ConfigOutput[T]) Run(s *T) error {
	if s == nil {
		panic("input is nil pointer")
	}

	c.Analyze(s)
	c.Traverse()
	return c.Execute()
}

type node struct {
	Children   []*node
	Value      reflect.Value
	FieldName  string
	TypeString string
	Parent     *node
}

func (n *node) String() string {
	buf := bytes.NewBuffer(nil)

	if len(n.Children) > 0 {
		fmt.Fprintf(buf, "STRUCT (%d):\n", len(n.Children))
		for _, n := range n.Children {
			s := n.String()
			lines := strings.Split(s, "\n")
			for _, line := range lines {
				fmt.Fprintf(buf, "  %s\n", line)

			}
		}
		fmt.Fprintln(buf, "END")
	} else {
		fmt.Fprintf(buf, "%s: %s (%v)", n.FieldName, n.TypeString, n.Value.Interface())
	}
	return buf.String()
}

func NewConfig[T any]() *ConfigOutput[T] {
	return &ConfigOutput[T]{}
}

func traverseTree(n *node) func() error {
	if len(n.Children) > 0 {
		return makeList(n)
	} else {
		return updateVal(n)
	}
}

func makeList(n *node) func() error {
	manager := makeManager(n)

	return func() error {
		for {
			end := manager.Tui()
			*manager = *makeManager(n)

			if end {
				break
			}
		}
		return nil
	}
}

func makeManager(n *node) *command.Manager {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})

	for i := 0; i < len(n.Children); i++ {
		child := n.Children[i]
		f := traverseTree(child)
		manager.Register(child.FieldName, fmt.Sprint(child.Value.Interface()), f)
	}

	return &manager
}

func updateVal(n *node) func() error {
	return func() error {
		var err error
		var vInt int
		var vFloat float64
		var vBool bool
		var vStr string

	inputVal:
		vStr = input.GetInput("New value")
		switch n.Value.Kind() {
		case reflect.Int:
			vInt, err = strconv.Atoi(vStr)
			if err != nil {
				fmt.Println("could not parse input. try again")
				goto inputVal
			}
			n.Value.Set(reflect.ValueOf(vInt))
		case reflect.Float64:
			vFloat, err = strconv.ParseFloat(vStr, 64)
			if err != nil {
				fmt.Println("could not parse input. try again")
				goto inputVal
			}
			n.Value.Set(reflect.ValueOf(vFloat))
		case reflect.Bool:
			vBool, err = strconv.ParseBool(vStr)
			if err != nil {
				fmt.Println("could not parse input. try again")
				goto inputVal
			}
			n.Value.Set(reflect.ValueOf(vBool))
		case reflect.String:
			n.Value.Set(reflect.ValueOf(vStr))
		default:
			return fmt.Errorf("invalid type: %s", n.Value.Kind())
		}
		return nil
	}
}

func makeTree(v reflect.Value) *node {
	children := []*node{}
	if v.Kind() == reflect.Pointer {
		return makeTree(v.Elem())
	}

	if isScalarVal(v) {
		return &node{Value: v}
	}

	if v.Kind() == reflect.Struct {
		parent := node{}
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			f := v.Field(i)
			n := makeTree(f)
			t := v.Type().Field(i)
			n.FieldName = t.Name
			n.TypeString = t.Type.Name()
			n.Parent = &parent
			children = append(children, n)
		}

		parent.Children = children

		return &parent
	}
	panic(fmt.Errorf("illegal type found: %s", v.Kind()))
}

func isScalarVal(v reflect.Value) bool {
	return isScalar(v.Kind())
}

func isScalar(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Float64, reflect.String, reflect.Bool:
		return true
	default:
		return false
	}
}

// produce an AST
// traverse the AST and produce a command tree
// allow scalar changes
// reflect changes to the object
