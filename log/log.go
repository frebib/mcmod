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

type contextKey string

const ContextKey contextKey = "log"

func FromContext(ctx context.Context) (log *logrus.Entry) {
	logObj := ctx.Value(ContextKey)
	log, ok := logObj.(*logrus.Entry)
	if !ok {
		log = logrus.WithContext(ctx)
	}
	return AddPrefix(log, 1)
}

// SetContextLogger updates the logger in the given context, or provides the
// default logger. It also provides a default context if none is given
func SetContextLogger(ctx context.Context, log *logrus.Entry) (context.Context, *logrus.Entry) {
	if ctx == nil {
		ctx = context.Background()
	}
	if log == nil {
		log = logrus.WithContext(ctx)
	}
	log = AddPrefix(log, 1)
	return context.WithValue(ctx, ContextKey, log), log
}

// AddPrefix uses the runtime information of the calling function to provide an
// intelligent approximation of what the log prefix should be. It sets the field
// "prefix" in the logger. The form <package>/<file> is used, where the
// <package> is the last element of the Golang package, and the <file> is the
// non-extension part of the filename, excluding the directory path. If the
// <package> and <file> are the same, no concatenation is used, and the single
// word is used instead. For example, the prefix for this file would be "log"
func AddPrefix(log *logrus.Entry, offset int) *logrus.Entry {
	// Source https://stackoverflow.com/a/25265493
	pc, filePath, _, ok := runtime.Caller(1 + offset)
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
