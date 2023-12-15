package pkgError

import (
	"fmt"
	"strings"
)

type packageError struct {
	err string
	pkg string
}

func (e *packageError) Error() string {
	return fmt.Sprintf("go-cli-tools/%s: %s", e.pkg, e.err)
}

func WrapError(pkg string, err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*packageError); ok {
		return e
	}

	str := err.Error()

	if strings.HasPrefix(str, "go-cli-tools") {
		return err
	}

	return Error(pkg, str)
}

func WrapErrorFactory(pkg string) func(error) error {
	return func(err error) error {
		return WrapError(pkg, err)
	}
}

func Error(pkg, err string) error {
	return &packageError{err, pkg}
}

func Errorf(pkg, err string, a ...any) error {
	formatted := fmt.Sprintf(err, a...)
	return &packageError{formatted, pkg}
}

func ErrorFactory(pkg string) func(string) error {
	return func(e string) error {
		return Error(pkg, e)
	}
}

func ErrorfFactory(pkg string) func(string, ...any) error {
	return func(e string, a ...any) error {
		return Errorf(pkg, e, a...)
	}
}
