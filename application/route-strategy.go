package application

import (
	"strings"
)

func WildcardStrategy(key string, ctx Context) Host {

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

func UrlStrategy(key string, ctx Context) Host {
	panic("Not implemented yet!")
}
