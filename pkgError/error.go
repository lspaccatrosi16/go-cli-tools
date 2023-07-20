package pkgError

import "fmt"

type packageError struct {
	err string
	pkg string
}

func (e *packageError) Error() string {
	return fmt.Sprintf("go-cli-tools/%s: %s", e.pkg, e.err)
}

func WrapError(pkg string, err error) error {
	str := err.Error()
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
