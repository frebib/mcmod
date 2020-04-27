package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/frebib/mcmod/api"
	modlog "github.com/frebib/mcmod/log"
	"github.com/frebib/mcmod/util"
	"github.com/urfave/cli/v2"
)

var (
	flagCount = cli.UintFlag{
		Name:    "count",
		Usage:   "number of results to display",
		Aliases: []string{"n"},
		Value:   10,
	}
	flagAll = cli.BoolFlag{
		Name: "all",
		Usage: fmt.Sprintf("show all results, not just the first %d",
			flagCount.Value),
		Aliases: []string{"a"},
		Value:   false,
	}
	Search = &cli.Command{
		Name:      "search",
		Usage:     "search for a mod",
		Action:    cmdDoModSearch,
		ArgsUsage: "<term>",
		Flags: []cli.Flag{
			&flagAll,
			&flagCount,
			&flagVersion,
		},
	}
)

func cmdDoModSearch(c *cli.Context) (err error) {
	ctx := c.Context
	log := modlog.FromContext(ctx)

	if c.NArg() < 1 {
		log.Error("missing required arg: " + c.Command.ArgsUsage)
		return cli.ShowSubcommandHelp(c)
	}

	term := strings.Join(c.Args().Slice(), " ")
	results, err := api.ClientFromContext(ctx).AddonSearch(ctx,
		api.AddonSearchOption{
			GameId:      api.GameMinecraft,
			Sort:        api.AddonSortPopularity,
			GameVersion: c.String(flagVersion.Name),
			Filter:      term,
		},
	)

	// TODO: Template output fields with text/template
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprint(w, "ID\tName\tDownloads\tLast Updated\tSlug\tVersions\n")

	// Attempt to do a better search, with a fuzzy search library
	for idx, mod := range results {
		versions := mod.SupportedVersions().LatestPatches().Strings()
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			mod.ID, util.EllipsiseString(mod.Name, 32),
			humanize.SIWithDigits(mod.DownloadCount, 2, ""),
			humanize.Time(mod.DateModified),
			mod.Slug,
			strings.Join(versions, ", "),
		)
		// Only show the max amount of results, if not displaying all
		if !c.Bool(flagAll.Name) && uint(idx) >= c.Uint(flagCount.Name)-1 {
			break
		}
	}
	return w.Flush()
}
