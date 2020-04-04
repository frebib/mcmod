package api

import (
	"context"
	"fmt"
	"strconv"

	modlog "github.com/frebib/mcmod/log"
)

func (c *ApiClient) Lookup(ctx context.Context, modNameId string) (addon *Addon, err error) {
	log := modlog.FromContext(ctx)

	// If input is not an int, assume it's a mod name
	if id, err := strconv.Atoi(modNameId); err != nil {
		log = log.WithField("name", modNameId)
		log.Debug("searching for mod by name")

		results, err := c.AddonSearch(ctx, AddonSearchOption{
			GameId: GameMinecraft,
			Filter: modNameId,
		})
		if err != nil {
			log.WithError(err).
				Errorf("failed to search")
		}

		addon = results.FindBySlug(modNameId)
		if addon != nil {
			log.Debug("matched mod by slug")
		} else {
			addon = results.FindByName(modNameId)
			if addon != nil {
				log.Debug("matched mod by fileName")
			}
		}
		if addon == nil {
			return nil, &ErrNoSuchAddon{Name: modNameId}
		}
	} else {
		log.WithField("id", id).
			Debug("fetching mod metadata by id")

		addon, err = c.AddonByID(ctx, id)
		if err != nil {
			log.WithError(err).
				Errorf("failed to query")
		}

		// Sanity-check that we have the right mod
		if addon == nil {
			return nil, &ErrNoSuchAddon{ID: id}
		} else if addon.ID != id {
			format := "mismatched local id %d to curseforge id %d. aborting"
			log.Warn(fmt.Sprintf(format, id, addon.ID))
			return nil, nil
		}
	}

	return addon, err
}
