package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/kubeshop/tracetest/server/analytics"
)

func jsonEscape(text any) string {
	initial, err := json.Marshal(text)
	if err != nil {
		panic(err)
	}

	encoded, err := json.Marshal(string(initial))
	if err != nil {
		panic(err)
	}

	formatted := string(encoded)
	return strings.Trim(formatted, `"`)
}

func spaHandler(prefix, staticPath, indexPath string, tplVars map[string]string) http.HandlerFunc {
	var fileMatcher = regexp.MustCompile(`\.[a-zA-Z]*$`)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			tpl, err := template.ParseFiles(filepath.Join(staticPath, indexPath))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if err = tpl.Execute(w, tplVars); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		} else {
			http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
		}
	}

	return http.StripPrefix(prefix, http.HandlerFunc(handler)).ServeHTTP
}

type spaConfig interface {
	ServerPathPrefix() string
	DemoEnabled() []string
	DemoEndpoints() map[string]string
	ExperimentalFeatures() []string
}

func SPAHandler(conf spaConfig, analyticsEnabled bool, serverID, version, env string) http.HandlerFunc {
	pathPrefix := conf.ServerPathPrefix()
	return spaHandler(
		pathPrefix,
		"./html",
		"index.html",
		map[string]string{
			"AnalyticsKey":         analytics.FrontendKey,
			"AnalyticsEnabled":     fmt.Sprintf("%t", analyticsEnabled),
			"ServerPathPrefix":     fmt.Sprintf("%s/", pathPrefix),
			"ServerID":             serverID,
			"AppVersion":           version,
			"Env":                  env,
			"DemoEnabled":          jsonEscape(conf.DemoEnabled()),
			"DemoEndpoints":        jsonEscape(conf.DemoEndpoints()),
			"ExperimentalFeatures": jsonEscape(conf.ExperimentalFeatures()),
		},
	)
}
