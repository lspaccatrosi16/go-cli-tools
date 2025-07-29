package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/lspaccatrosi16/go-cli-tools/internal/pkgError"
)

var parsed bool

var args []string
var version string

var wrap = pkgError.WrapErrorFactory("args")

var writer = os.Stdout

func SetVersion(v string) {
	version = v
}

var SetWriter = func(w *os.File) {
	writer = w
}

func ParseOpts() error {
	if !flag.Parsed() {
		transformEntries()
		flag.Parse()
		rem := flag.Args()
		if h, err := GetFlagValue[bool]("help"); err == nil && h {
			fmt.Fprintln(writer, Usage())
			os.Exit(0)
		} else if err != nil {
			return err
		}

		if v, err := GetFlagValue[bool]("version"); err == nil && v {
			fmt.Fprintln(writer, version)
			os.Exit(0)
		} else if err != nil {
			return err
		}

		args = rem
		parsed = true
	}
	return nil
}

func GetArgs() []string {
	if !parsed {
		return nil
	}
	return args
}
