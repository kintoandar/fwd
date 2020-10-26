package version

import (
	"fmt"
	"runtime"
)

// GitCommit filled in by the compiler.
var GitCommit string

// Version number.
const Version = "1.1.0"

// GoVersion returns the version of the go runtime used to compile the binary.
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary.
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
