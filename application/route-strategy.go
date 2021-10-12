package application

import (
	"github.com/nekinci/paas/core"
	"strings"
)

func WildcardStrategy(key string, ctx Context) core.Host {
	if key == "" {
		return nil
	}

	splittedKey := strings.Split(key, ".")
	newKey := splittedKey[0]

	if isReservedName(newKey) {
		return ctx.embeddedApplications[newKey]
	}

	app := ctx.validApplications[newKey]
	if app == nil {
		return nil
	}

	a := *app
	if a.GetStatus() == WAITING {
		return nil
	}

	return *ctx.validApplications[newKey]
}

func UrlStrategy(key string, ctx Context) core.Host {
	// Not implemented yet.
	return nil
}
