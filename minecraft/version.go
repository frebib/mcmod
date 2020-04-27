package minecraft

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var verRegex = regexp.MustCompile(`(?i)^(?:(?P<release>alpha|beta|indev|infdev) )?v?` +
	`(?P<major>\d)\.(?P<minor>\d{1,2})(?:\.(?P<patch>\d{1,2}))?(?:_(?P<build>\d{2}))?$`)

type Phase string

const (
	Release = Phase("release")
	Alpha   = Phase("alpha")
	Beta    = Phase("beta")
	Indev   = Phase("indev")
	Infdev  = Phase("infdev")
)

func ParsePhase(s string) Phase {
	s = strings.ToLower(s)
	switch Phase(s) {
	case Alpha:
		return Alpha
	case Beta:
		return Beta
	case Indev:
		return Indev
	case Infdev:
		return Infdev
	default:
		return Release
	}
}

type Version struct {
	Major   int
	Minor   int
	Patch   int
	Build   int
	Release Phase
}

func (v Version) String() string {
	var s string
	if v.Release != "" && v.Release != Release {
		s = strings.ToUpper(fmt.Sprintf("%c", v.Release[0])) + string(v.Release)[1:] + " "
	}
	s += fmt.Sprintf("%d", v.Major)
	if v.Minor >= 0 {
		s += fmt.Sprintf(".%d", v.Minor)
	}
	if v.Patch >= 0 {
		s += fmt.Sprintf(".%d", v.Patch)
	}
	if v.Build >= 0 {
		s += fmt.Sprintf("_%02d", v.Build)
	}
	return s
}

type Versions []Version

func (vs Versions) Len() int {
	return len(vs)
}

func (vs Versions) Less(i, j int) bool {
	return vs[i].LessThan(vs[j])
}

func (vs Versions) Swap(i, j int) {
	vs[j], vs[i] = vs[i], vs[j]
}

func (vs Versions) Strings() []string {
	var strs = make([]string, len(vs))
	for idx, ver := range vs {
		strs[idx] = ver.String()
	}
	return strs
}

func (vs Versions) LatestPatches() *Versions {
	latestFor := make(map[string]Version, 0)
	for _, ver := range vs {
		sigStr := fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
		// If there is already a version for the significant, and it is greater
		// don't promote the lower version and continue to next
		if latest, ok := latestFor[sigStr]; ok &&
			latest.LessThan(ver) {
			continue
		}
		latestFor[sigStr] = ver
	}

	versions := make(Versions, 0)
	for _, ver := range latestFor {
		versions = append(versions, ver)
	}
	sort.Sort(versions)
	return &versions
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", v.String())), nil
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	return ParseInto(v, str)
}

func (v Version) Compare(o Version) int {
	cmp := o.Major - v.Major
	if cmp != 0 {
		return cmp
	}
	cmp = o.Minor - v.Minor
	if cmp != 0 {
		return cmp
	}
	cmp = o.Patch - v.Patch
	if cmp != 0 {
		return cmp
	}
	return o.Build - v.Build
}

func (v Version) LessThan(o Version) bool {
	return v.Compare(o) < 1
}

func (v Version) GreaterThan(o Version) bool {
	return v.Compare(o) > 1
}

func MustParse(s string) *Version {
	ver, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return ver
}

func Parse(s string) (*Version, error) {
	version := new(Version)
	err := ParseInto(version, s)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func ParseInto(v *Version, s string) error {
	parts := verRegex.FindStringSubmatch(s)
	if len(parts) < 1 || parts[0] == "" {
		return ErrInvalidVersion{s}
	}
	v.Minor = -1
	v.Patch = -1
	v.Build = -1
	v.Release = Release
	for idx, name := range verRegex.SubexpNames() {
		// Skip the first value, it's the whole match in one
		if name == "" || idx == 0 {
			continue
		}
		value := parts[idx]
		var err error
		switch name {
		case "release":
			v.Release = ParsePhase(value)
		case "major":
			v.Major, err = strconv.Atoi(value)
		case "minor":
			v.Minor, err = strconv.Atoi(value)
		case "patch":
			if value != "" {
				v.Patch, err = strconv.Atoi(value)
			}
		case "build":
			if value != "" {
				v.Build, err = strconv.Atoi(value)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
