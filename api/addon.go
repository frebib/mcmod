package api

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type Addon struct {
	ID                     int                          `json:"id"`
	Name                   string                       `json:"name"`
	Authors                []AddonAuthor                `json:"authors"`
	Attachments            []AddonAttachment            `json:"attachments"`
	WebsiteURL             string                       `json:"websiteUrl"`
	GameID                 int                          `json:"gameId"`
	Summary                string                       `json:"summary"`
	DefaultFileID          int                          `json:"defaultFileId"`
	DownloadCount          float64                      `json:"downloadCount"`
	LatestFiles            []AddonLatestFile            `json:"latestFiles"`
	Categories             []AddonCategory              `json:"categories"`
	Status                 int                          `json:"status"`
	PrimaryCategoryID      int                          `json:"primaryCategoryId"`
	CategorySection        AddonCategorySection         `json:"categorySection"`
	Slug                   string                       `json:"slug"`
	GameVersionLatestFiles []AddonGameVersionLatestFile `json:"gameVersionLatestFiles"`
	IsFeatured             bool                         `json:"isFeatured"`
	PopularityScore        float64                      `json:"popularityScore"`
	GamePopularityRank     int                          `json:"gamePopularityRank"`
	PrimaryLanguage        string                       `json:"primaryLanguage"`
	GameSlug               string                       `json:"gameSlug"`
	GameName               string                       `json:"gameName"`
	PortalName             string                       `json:"portalName"`
	DateModified           time.Time                    `json:"dateModified"`
	DateCreated            time.Time                    `json:"dateCreated"`
	DateReleased           time.Time                    `json:"dateReleased"`
	IsAvailable            bool                         `json:"isAvailable"`
	IsExperiemental        bool                         `json:"isExperiemental"`
}
type AddonAuthor struct {
	Name              string      `json:"name"`
	URL               string      `json:"url"`
	ProjectID         int         `json:"projectId"`
	ID                int         `json:"id"`
	ProjectTitleID    interface{} `json:"projectTitleId"`
	ProjectTitleTitle interface{} `json:"projectTitleTitle"`
	UserID            int         `json:"userId"`
	TwitchID          int         `json:"twitchId"`
}
type AddonAttachment struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"projectId"`
	Description  string `json:"description"`
	IsDefault    bool   `json:"isDefault"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Status       int    `json:"status"`
}
type AddonModule struct {
	Foldername  string `json:"foldername"`
	Fingerprint int64  `json:"fingerprint"`
	Type        int    `json:"type"`
}
type AddonSortableGameVersion struct {
	GameVersionPadded      string    `json:"gameVersionPadded"`
	GameVersion            string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionName        string    `json:"gameVersionName"`
}
type AddonLatestFile struct {
	ID                         int                        `json:"id"`
	DisplayName                string                     `json:"displayName"`
	FileName                   string                     `json:"fileName"`
	FileDate                   time.Time                  `json:"fileDate"`
	FileLength                 int                        `json:"fileLength"`
	ReleaseType                int                        `json:"releaseType"`
	FileStatus                 int                        `json:"fileStatus"`
	DownloadURL                string                     `json:"downloadUrl"`
	IsAlternate                bool                       `json:"isAlternate"`
	AlternateFileID            int                        `json:"alternateFileId"`
	Dependencies               []interface{}              `json:"dependencies"`
	IsAvailable                bool                       `json:"isAvailable"`
	Modules                    []AddonModule              `json:"modules"`
	PackageFingerprint         int64                      `json:"packageFingerprint"`
	GameVersion                []string                   `json:"gameVersion"`
	SortableGameVersion        []AddonSortableGameVersion `json:"sortableGameVersion"`
	InstallMetadata            interface{}                `json:"installMetadata"`
	Changelog                  interface{}                `json:"changelog"`
	HasInstallScript           bool                       `json:"hasInstallScript"`
	IsCompatibleWithClient     bool                       `json:"isCompatibleWithClient"`
	CategorySectionPackageType int                        `json:"categorySectionPackageType"`
	RestrictProjectFileAccess  int                        `json:"restrictProjectFileAccess"`
	ProjectStatus              int                        `json:"projectStatus"`
	RenderCacheID              int                        `json:"renderCacheId"`
	FileLegacyMappingID        interface{}                `json:"fileLegacyMappingId"`
	ProjectID                  int                        `json:"projectId"`
	ParentProjectFileID        interface{}                `json:"parentProjectFileId"`
	ParentFileLegacyMappingID  interface{}                `json:"parentFileLegacyMappingId"`
	FileTypeID                 interface{}                `json:"fileTypeId"`
	ExposeAsAlternative        interface{}                `json:"exposeAsAlternative"`
	PackageFingerprintID       int                        `json:"packageFingerprintId"`
	GameVersionDateReleased    time.Time                  `json:"gameVersionDateReleased"`
	GameVersionMappingID       int                        `json:"gameVersionMappingId"`
	GameVersionID              int                        `json:"gameVersionId"`
	GameID                     int                        `json:"gameId"`
	IsServerPack               bool                       `json:"isServerPack"`
	ServerPackFileID           interface{}                `json:"serverPackFileId"`
	GameVersionFlavor          interface{}                `json:"gameVersionFlavor"`
}
type AddonCategory struct {
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	AvatarURL  string `json:"avatarUrl"`
	ParentID   int    `json:"parentId"`
	RootID     int    `json:"rootId"`
	ProjectID  int    `json:"projectId"`
	AvatarID   int    `json:"avatarId"`
	GameID     int    `json:"gameId"`
}
type AddonCategorySection struct {
	ID                      int         `json:"id"`
	GameID                  int         `json:"gameId"`
	Name                    string      `json:"name"`
	PackageType             int         `json:"packageType"`
	Path                    string      `json:"path"`
	InitialInclusionPattern string      `json:"initialInclusionPattern"`
	ExtraIncludePattern     interface{} `json:"extraIncludePattern"`
	GameCategoryID          int         `json:"gameCategoryId"`
}
type AddonGameVersionLatestFile struct {
	GameVersion       string      `json:"gameVersion"`
	ProjectFileID     int         `json:"projectFileId"`
	ProjectFileName   string      `json:"projectFileName"`
	FileType          int         `json:"fileType"`
	GameVersionFlavor interface{} `json:"gameVersionFlavor"`
}

type SearchResult []Addon

func (sr *SearchResult) FindBySlug(slug string) *Addon {
	slug = strings.ToLower(slug)
	for idx, addon := range *sr {
		if strings.ToLower(addon.Slug) == slug {
			return &(*sr)[idx]
		}
	}
	return nil
}

func (sr *SearchResult) FindByName(name string) *Addon {
	name = strings.ToLower(name)
	for idx, addon := range *sr {
		if strings.ToLower(addon.Name) == name {
			return &(*sr)[idx]
		}
	}
	return nil
}

type AddonSearchOption struct {
	CategoryID  int    `url:"categoryID,omitempty"`
	SectionId   int    `url:"sectionId,omitempty"`
	GameId      int    `url:"gameId"`
	GameVersion string `url:"gameVersion,omitempty"`
	Index       int    `url:"index,omitempty"`
	PageSize    int    `url:"pageSize"`
	Filter      string `url:"searchFilter"`
	Slug        string `url:"slug,omitempty"`
	Sort        int    `url:"sort,omitempty"`
}

func setDefaultUnsetOptions(opts *AddonSearchOption) *AddonSearchOption {
	if opts.PageSize == 0 {
		opts.PageSize = 9999
	}
	return opts
}

func (c *ApiClient) AddonSearch(ctx context.Context, opts AddonSearchOption) (*SearchResult, error) {
	params, err := query.Values(setDefaultUnsetOptions(&opts))
	if err != nil {
		return nil, err
	}
	queryUrl, err := buildURLParams(c.ApiUrl, "v2/addon/search", &params)
	if err != nil {
		return nil, err
	}

	resp, err := fetchJSON(ctx, c.HttpClient, "GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result SearchResult
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

func (c *ApiClient) AddonByID(ctx context.Context, id int) (*Addon, error) {
	queryUrl, err := buildURL(c.ApiUrl, "v2/addon/"+strconv.Itoa(id), "")
	if err != nil {
		return nil, err
	}

	resp, err := fetchJSON(ctx, c.HttpClient, "GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var addon Addon
	return &addon, json.NewDecoder(resp.Body).Decode(&addon)
}
