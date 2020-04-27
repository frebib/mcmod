package api

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/frebib/mcmod/util"
)

type Files []File

type File struct {
	ID                      int          `json:"id"`
	DisplayName             string       `json:"displayName"`
	FileName                string       `json:"fileName"`
	FileDate                time.Time    `json:"fileDate"`
	FileLength              int          `json:"fileLength"`
	ReleaseType             ReleaseType  `json:"releaseType"`
	FileStatus              int          `json:"fileStatus"`
	DownloadURL             string       `json:"downloadUrl"`
	IsAlternate             bool         `json:"isAlternate"`
	AlternateFileID         int          `json:"alternateFileId"`
	Dependencies            []Dependency `json:"dependencies"`
	IsAvailable             bool         `json:"isAvailable"`
	Modules                 []Module     `json:"modules"`
	PackageFingerprint      int64        `json:"packageFingerprint"`
	GameVersion             []string     `json:"gameVersion"`
	InstallMetadata         interface{}  `json:"installMetadata"`
	ServerPackFileID        interface{}  `json:"serverPackFileId"`
	HasInstallScript        bool         `json:"hasInstallScript"`
	GameVersionDateReleased time.Time    `json:"gameVersionDateReleased"`
	GameVersionFlavor       interface{}  `json:"gameVersionFlavor"`
}

type Dependency struct {
	ID      int `json:"id"`
	AddonID int `json:"addonId"`
	Type    int `json:"type"`
	FileID  int `json:"fileId"`
}

type Module struct {
	Foldername  string `json:"foldername"`
	Fingerprint int64  `json:"fingerprint"`
}

type ReleaseType int

const (
	ReleaseAny     ReleaseType = -1
	ReleaseUnknown ReleaseType = 0
	ReleaseRelease ReleaseType = 1
	ReleaseBeta    ReleaseType = 2
	ReleaseAlpha   ReleaseType = 3
)

func (f ReleaseType) String() string {
	return [...]string{"any", "unknown", "release", "beta", "alpha"}[f+1]
}
func ParseReleaseType(s string) ReleaseType {
	switch strings.ToLower(s) {
	case "any":
		return ReleaseAny
	case "release":
		return ReleaseRelease
	case "beta":
		return ReleaseBeta
	case "alpha":
		return ReleaseAlpha
	}
	return ReleaseUnknown
}

func (fs Files) Len() int {
	return len(fs)
}

func (fs Files) Less(i, j int) bool {
	return fs[i].FileDate.UnixNano() > fs[j].FileDate.UnixNano()
}

func (fs Files) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

var _ sort.Interface = &Files{}

func (c *ApiClient) Files(ctx context.Context, mod int) (files Files, err error) {
	path := fmt.Sprintf("v2/addon/%d/files", mod)
	queryUrl, err := buildURL(c.ApiUrl, path, "")
	if err != nil {
		return nil, err
	}

	_, err = fetchJSON(ctx, c.HttpClient, "GET", queryUrl, nil, &files)
	return
}

func (c *ApiClient) FileByID(ctx context.Context, addon, file int) (res *File, err error) {
	path := fmt.Sprintf("v2/addon/%d/file/%d/download-url", addon, file)
	queryUrl, err := buildURL(c.ApiUrl, path, "")
	if err != nil {
		return nil, err
	}

	_, err = fetchJSON(ctx, c.HttpClient, "GET", queryUrl, nil, &res)
	return
}

type FileFilter struct {
	FilterFunc func(*File) bool
	AfterFunc  func(Files) error
}

func FileFilterVersion(ver string) FileFilter {
	var verClos = ver
	return FileFilter{
		func(file *File) bool {
			return util.StringInSlice(file.GameVersion, verClos)
		},
		nil,
	}
}

func FileFilterRelease(release ReleaseType) FileFilter {
	var releaseClos = release
	return FileFilter{
		func(file *File) bool {
			return file.ReleaseType <= releaseClos
		},
		nil,
	}
}

func (fs *Files) Filter(filters []FileFilter) (Files, error) {
	// Copy the input slice, instead of mutating it
	var all Files = append(make(Files, 0), *fs...)

	// Apply each filter one at a time
	for _, filter := range filters {
		var remain Files = make(Files, 0)
		for _, file := range all {
			if filter.FilterFunc(&file) {
				remain = append(remain, file)
			}
		}
		if filter.AfterFunc != nil {
			err := filter.AfterFunc(remain)
			if err != nil {
				return remain, err
			}
		}

		all = remain
		if len(all) < 1 {
			break
		}
	}
	return all, nil
}
