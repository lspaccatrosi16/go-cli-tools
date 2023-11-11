package credential

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/config"
	"github.com/lspaccatrosi16/go-cli-tools/gbin"
	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type wrappedCredential struct {
	Name        string
	Description string
	Cred        Credential
}

type credentialmanager struct {
	Credentials []*wrappedCredential
}

func (c *credentialmanager) pick(appName string, loop bool) (*Credential, error) {
start:
	var chosenCredential Credential
	var opts []input.SelectOption
	var defaultCredentials *Credential
	var err error

	opts = append(opts, input.SelectOption{Name: "Back", Value: "e"})

	if appName != "" {
		defaultCredentials, err = getLegacyCredentials(appName)
		if err != nil {
			return nil, err
		}

		if defaultCredentials.Key != "" && defaultCredentials.Secret != "" {
			opts = append(opts, input.SelectOption{Name: "Use default application credentials", Value: "d"})
		}
	}

	opts = append(opts, input.SelectOption{Name: "Use a saved credential", Value: "s"})
	opts = append(opts, input.SelectOption{Name: "Add a new credential", Value: "n"})
	opts = append(opts, input.SelectOption{Name: "Remove a saved credential", Value: "r"})

	selected, err := input.GetSelection("Select Preference", opts)

	if err != nil {
		return nil, wrap(err)
	}

	switch selected {
	case "e":
		return nil, nil
	case "d":
		return defaultCredentials, nil
	case "s":
		manager := command.NewManager(command.ManagerConfig{Searchable: true})

		for _, wc := range c.Credentials {
			manager.Register(wc.Name, wc.Description, setPtr(&chosenCredential, wc.Cred))
		}

		back := manager.Tui()
		if back {
			goto start
		}

	case "n":
		created := newCredential()
		c.Credentials = append(c.Credentials, created)
		chosenCredential = created.Cred
	case "r":
		manager := command.NewManager(command.ManagerConfig{Searchable: true})
		for _, wc := range c.Credentials {
			manager.Register(wc.Name, wc.Description, removeCredential(c, wc))
		}

		manager.Tui()
		goto start

	}

	makeDefault, err := input.GetConfirmSelection("Make this the default credential for this app")

	if err != nil {
		return nil, wrap(err)
	}

	if makeDefault {
		err = setLegacyCredentials(appName, &chosenCredential)
		if err != nil {
			return nil, err
		}
	}

	if loop {
		goto start
	}

	return &chosenCredential, nil
}

func newCredential() *wrappedCredential {
	fmt.Println("New Credential")
	name := input.GetInput("Name")
	description := input.GetInput("Description")
	key := input.GetInput("Key")
	secret := input.GetInput("Secret")
	return &wrappedCredential{
		Name:        name,
		Description: description,
		Cred: Credential{
			Key:    key,
			Secret: secret,
		},
	}
}

func removeCredential(manager *credentialmanager, cred *wrappedCredential) func() error {
	return func() error {
		newCreds := []*wrappedCredential{}

		for _, wc := range manager.Credentials {
			if wc != cred {
				newCreds = append(newCreds, wc)
			}
		}
		if len(newCreds)+1 != len(manager.Credentials) {
			return fmt.Errorf("new list of length %d is not of expected length %d", len(newCreds), len(manager.Credentials)-1)
		}

		manager.Credentials = newCreds

		return nil
	}
}

func setPtr(p *Credential, v Credential) func() error {
	return func() error {
		*p = v
		return nil
	}
}

func centralCredLocation() (string, error) {
	cpath, err := config.GetConfigPath("gct-credmanager")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", wrap(err)
	}

	return filepath.Join(cpath, "credstore"), nil
}

func saveManager(manager *credentialmanager) error {
	path, err := centralCredLocation()
	if err != nil {
		return err
	}

	enc := gbin.NewEncoder[credentialmanager]()

	src, err := enc.EncodeStream(manager)
	if err != nil {
		return wrap(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return wrap(err)
	}
	defer f.Close()

	io.Copy(f, src)
	return nil
}

func loadManager() (*credentialmanager, error) {
	path, err := centralCredLocation()
	if err != nil {
		return nil, err
	}

	dec := gbin.NewDecoder[credentialmanager]()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &credentialmanager{}, nil
		}
		return nil, wrap(err)
	}

	defer f.Close()

	manager, err := dec.DecodeStream(f)
	if err != nil {
		return nil, wrap(err)
	}

	return manager, nil
}
