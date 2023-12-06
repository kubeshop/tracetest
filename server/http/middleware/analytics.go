package middleware

import (
	"context"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/analytics"
)

var (
	matchFirstCap     = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap       = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchResourceName = regexp.MustCompile(`(\w)(\.)(\w)`)
)

const EventPropertiesKey key = "x-event-properties"

func AnalyticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeName := mux.CurrentRoute(r).GetName()
		machineID := r.Header.Get("x-client-id")
		source := r.Header.Get("x-source")
		eventProperties := r.Header.Get("x-event-properties")

		eventData := map[string]string{
			"source": source,
		}
		eventData = analytics.InjectProperties(eventData, eventProperties)

		if routeName != "" {
			analytics.SendEvent(toWords(routeName), "test", machineID, &eventData)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, EventPropertiesKey, eventProperties)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func toWords(str string) string {
	if matchResourceName.MatchString(str) {
		return str
	}
	words := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	words = matchAllCap.ReplaceAllString(words, "${1} ${2}")
	return words
}

func EventPropertiesFromContext(ctx context.Context) string {
	eventProperties := ctx.Value(EventPropertiesKey)

	if eventProperties == nil {
		return ""
	}

	return eventProperties.(string)
}
