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
	"github.com/kubeshop/tracetest/server/config"
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

func SPAHandler(conf config.Config, serverID, version, env string) http.HandlerFunc {
	return spaHandler(
		conf.Server.PathPrefix,
		"./html",
		"index.html",
		map[string]string{
			"AnalyticsKey":         analytics.FrontendKey,
			"AnalyticsEnabled":     fmt.Sprintf("%t", conf.GA.Enabled),
			"ServerPathPrefix":     fmt.Sprintf("%s/", conf.Server.PathPrefix),
			"ServerID":             serverID,
			"AppVersion":           version,
			"Env":                  env,
			"DemoEnabled":          jsonEscape(conf.Demo.Enabled),
			"DemoEndpoints":        jsonEscape(conf.Demo.Endpoints),
			"ExperimentalFeatures": jsonEscape(conf.ExperimentalFeatures),
		},
	)
}
