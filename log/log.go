package log

import (
	"context"
	"regexp"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	versionRegex = regexp.MustCompile("^v\\d$")
)

func FromContext(ctx context.Context) (log *logrus.Entry) {
	logObj := ctx.Value("log")
	log, ok := logObj.(*logrus.Entry)
	if !ok {
		log = logrus.WithContext(ctx)
	}

	// Source https://stackoverflow.com/a/25265493
	pc, filePath, _, ok := runtime.Caller(2)
	if ok {
		// From the file name we're interested in the name before the dot
		// e.g.   src/thing/another/log.go  ->  log
		fileParts := strings.Split(filePath, "/")
		fileName := strings.Split(fileParts[len(fileParts)-1], ".")[0]
		pkgName := normalisePkgName(runtime.FuncForPC(pc).Name())

		if fileName != pkgName {
			pkgName = pkgName + "/" + fileName
		}

		log = log.WithField("prefix", pkgName)
	}

	return log
}

func normalisePkgName(s string) (pkg string) {
	// Split to remove the function name (and brackets)
	parts := strings.Split(s, ".")
	partIdx := len(parts) - 2

	// Cut out brackets from package names like
	//   github.com/urfave/cli/v2.(*Command).Run
	if parts[partIdx][0] == '(' {
		partIdx--
	}

	// Get the last part of the package name
	// e.g. github.com/frebib/mcmod   ->  mcmod
	// or   github.com/urfave/cli/v2  ->  cli (make sure to strip the version)
	parts = strings.Split(parts[partIdx], "/")
	pkg = parts[len(parts)-1]

	// Ensure to strip out the
	if versionRegex.Match([]byte(pkg)) {
		// we want the part before the `v2` or whatever
		return parts[len(parts)-2]
	}
	return pkg
}
