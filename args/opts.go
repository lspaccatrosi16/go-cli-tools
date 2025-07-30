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

func Version() {
	fmt.Fprintln(writer, version)
}

var cust_usage = func() {
	fmt.Fprintln(writer, Usage())
}

func UseCustomUsage(f func()) {
	cust_usage = f
}

func ParseOpts() error {
	if !flag.Parsed() {
		flag.Usage = cust_usage
		transformEntries()
		flag.Parse()
		rem := flag.Args()
		if h, err := GetFlagValue[bool]("help"); err == nil && h {
			cust_usage()
			os.Exit(0)
		}

		if v, err := GetFlagValue[bool]("version"); err == nil && v {
			Version()
			os.Exit(0)
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
