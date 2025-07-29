package args

import (
	"bytes"
	"flag"
	"fmt"
	"slices"
	"text/tabwriter"
)

type entry struct {
	name        string
	short       string
	description string
}

func (e *entry) Name() string {
	return e.name
}

func (e *entry) Short() string {
	return e.short
}

func (e *entry) Description() string {
	return e.description
}

type stringEntry struct {
	entry
	value string
}

func (e *stringEntry) Default() string {
	return e.value
}

type numberEntry struct {
	entry
	value int
}

func (e *numberEntry) Default() string {
	return fmt.Sprintf("%d", e.value)
}

type boolEntry struct {
	entry
	value bool
}

func (e *boolEntry) Default() string {
	return fmt.Sprintf("%t", e.value)
}

func NewStringEntry(name, short, description string, defaultVal string) ListEntry {
	return &stringEntry{
		entry: entry{
			name:        name,
			short:       short,
			description: description,
		},
		value: defaultVal,
	}
}

func NewNumberEntry(name, short, description string, defaultVal int) ListEntry {
	return &numberEntry{
		entry: entry{
			name:        name,
			short:       short,
			description: description,
		},
		value: defaultVal,
	}
}

func NewBoolEntry(name, short, description string, defaultVal bool) ListEntry {
	return &boolEntry{
		entry: entry{
			name:        name,
			short:       short,
			description: description,
		},
		value: defaultVal,
	}
}

type ListEntry interface {
	Name() string
	Short() string
	Description() string
	Default() string
}

var entries = map[string]ListEntry{
	"h": NewBoolEntry("help", "h", "show help", false),
	"v": NewBoolEntry("version", "v", "show version", false),
}

func UseCustomHelpTrigger() {
	delete(entries, "h")
}

func RegisterEntry(e ListEntry) {
	entries[e.Name()] = e
}

func Usage() string {
	buf := bytes.NewBuffer(nil)
	tr := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	ents := []ListEntry{}
	for _, v := range entries {
		ents = append(ents, v)
	}

	slices.SortFunc(ents, func(a, b ListEntry) int {
		if a.Name() < b.Name() {
			return -1
		}
		return 1
	})

	for _, v := range ents {
		fmt.Fprintf(tr, "%s\t-%s\t%s\t(%s)\n", v.Name(), v.Short(), v.Description(), v.Default())
	}
	tr.Flush()
	return buf.String()
}

var transformed map[string]any

func transformEntries() {
	ents := map[string]any{}
	for _, v := range entries {
		switch v := v.(type) {
		case *stringEntry:
			fv := flag.String(v.short, v.value, v.description)
			ents[v.name] = fv
		case *numberEntry:
			fv := flag.Int(v.short, v.value, v.description)
			ents[v.name] = fv
		case *boolEntry:
			fv := flag.Bool(v.short, v.value, v.description)
			ents[v.name] = fv
		}
	}
	transformed = ents
}

func GetFlagValue[T any](name string) (T, error) {
	var blank T
	if transformed == nil {
		return blank, wrap(fmt.Errorf("get flag value called before flag parse"))
	}

	v, ok := transformed[name]
	if !ok {
		return blank, wrap(fmt.Errorf("flag %s not found", name))
	}

	if fv, ok := v.(*T); ok {
		return *fv, nil
	}

	return blank, fmt.Errorf("flag %s is not of type %T", name, blank)
}
