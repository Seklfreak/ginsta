package ginsta

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func (g *Ginsta) gis(ctx context.Context, data string) (string, error) {
	var err error
	if g.rhxgis == "" {
		g.rhxgis, err = g.getRHXGIS(ctx)
		if err != nil {
			return "", err
		}
	}

	return md5Hash([]byte(g.rhxgis + ":" + data)), nil
}

func (g *Ginsta) getRHXGIS(ctx context.Context) (string, error) {
	profile, err := g.userRawProfileByUsername(ctx, "gfriendofficial")
	if err != nil {
		return "", err
	}

	if profile.RHXGIS == "" {
		return "", errors.New("unable to extract RHX GIS")
	}

	return profile.RHXGIS, nil
}

func md5Hash(text []byte) string {
	hasher := md5.New()
	hasher.Write(text)
	return hex.EncodeToString(hasher.Sum(nil))
}
