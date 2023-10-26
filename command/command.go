package command

import (
	"bytes"
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type cmd struct {
	Name        string
	Description string
	Exec        *func() error
}

func (c *cmd) Run() error {
	if c.Exec == nil {
		panic("exec property of command must not be nil")
	}
	return (*c.Exec)()
}

type datacmd struct {
	Name        string
	Description string
	Exec        *func() (any, error)
}

func (d *datacmd) Run() (any, error) {
	if d.Exec == nil {
		panic("exec property of datacommand must not be nil")
	}
	return (*d.Exec)()
}

type ManagerConfig struct {
	Searchable bool
}

type manager struct {
	cmds     []*cmd
	datacmds []*datacmd
	config   ManagerConfig
}

func (m *manager) Help() {
	maxCmdLength := 0
	cmds := []string{}
	descriptions := []string{}
	for _, cmd := range m.cmds {
		if len(cmd.Name) > maxCmdLength {
			maxCmdLength = len(cmd.Name)
		}
		cmds = append(cmds, cmd.Name)
		descriptions = append(descriptions, cmd.Description)
	}
	buf := bytes.NewBuffer(nil)

	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		desc := descriptions[i]

		fmt.Fprintf(buf, "%-*s : %s\n", maxCmdLength, cmd, desc)
	}

	fmt.Println(buf.String())
}

func (m *manager) Register(name string, description string, exec func() error) {
	newcmd := cmd{
		Name:        name,
		Description: description,
		Exec:        &exec,
	}
	m.cmds = append(m.cmds, &newcmd)
}

func (m *manager) RegisterData(name string, description string, exec func() (any, error)) {
	newcmd := datacmd{
		Name:        name,
		Description: description,
		Exec:        &exec,
	}
	m.datacmds = append(m.datacmds, &newcmd)
}

func (m *manager) Run(str string) {
	for _, cmd := range m.cmds {
		if cmd.Name == str {
			err := cmd.Run()
			if err != nil {
				fmt.Println("an error was encountered whilst running the command:\n", err.Error())
			}
			return
		}
	}

	fmt.Printf("command \"%s\" was not found\n", str)
	m.Help()
}

func (m *manager) RunData(str string) (any, error) {
	for _, cmd := range m.datacmds {
		if cmd.Name == str {
			data, err := cmd.Run()
			if err != nil {
				fmt.Println("an error was encountered whilst running the datacommand:\n", err.Error())
				return nil, err
			} else {
				return data, nil
			}
		}
	}

	estr := fmt.Sprintf("command \"%s\" was not found", str)
	fmt.Println(estr)

	m.Help()

	return nil, fmt.Errorf(estr)
}

func (m *manager) DataTui() (any, error) {
	maxCmdLen := 0
	names := []string{}
	descriptions := []string{}

	for _, cmd := range m.datacmds {
		if len(cmd.Name) > maxCmdLen {
			maxCmdLen = len(cmd.Name)
		}
		names = append(names, cmd.Name)
		descriptions = append(descriptions, cmd.Description)
	}

	selected := m.runTui(names, descriptions, maxCmdLen)

	if selected == "exit" {
		return nil, fmt.Errorf("no value selected")
	}

	return m.RunData(selected)
}

func (m *manager) Tui() bool {
	maxCmdLen := 0
	names := []string{}
	descriptions := []string{}

	for _, cmd := range m.cmds {
		if len(cmd.Name) > maxCmdLen {
			maxCmdLen = len(cmd.Name)
		}
		names = append(names, cmd.Name)
		descriptions = append(descriptions, cmd.Description)
	}

	selected := m.runTui(names, descriptions, maxCmdLen)

	if selected == "exit" {
		return true
	}

	m.Run(selected)
	return false
}

func (m *manager) runTui(names []string, descriptions []string, maxCmdLen int) string {
	options := []input.SelectOption{}

	options = append(options, input.SelectOption{Name: "Back", Value: "exit"})

	for i := 0; i < len(names); i++ {
		name := names[i]
		description := descriptions[i]
		options = append(options, input.SelectOption{Name: fmt.Sprintf("%-*s : %s", maxCmdLen, name, description), Value: name})

	}

	var selected string
	var err error

	if m.config.Searchable {
		selected, err = input.GetSearchableSelection("Select the command to execute", options)
	} else {
		selected, err = input.GetSelection("Select the command to execute", options)
	}

	if err != nil {
		panic(err)
	}

	return selected
}

func NewManager(config ManagerConfig) manager {
	return manager{config: config}
}
