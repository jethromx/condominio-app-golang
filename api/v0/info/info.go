package info

import (
	"sync"
)

type ReleaseInfo struct {
	Date    string `json:"date"`
	Commit  string `json:"commit"`
	Name    string `json:"name"`
	Uuid    string `json:"uuid"`
	Version string `json:"version"`
}

type AppInfo struct {
	ServiceBaseURL string       `json:"_serviceBaseURL"`
	ServiceVersion string       `json:"_serviceVersion"`
	ApiBaseURL     string       `json:"apiBaseURL"`
	ApiVersion     string       `json:"apiVersion"`
	ServiceLocator string       `json:"serviceLocator"`
	ReleaseInfo    *ReleaseInfo `json:"releaseInfo"`
}

var (
	instance *AppInfo
	once     sync.Once
)

func getAppInfo() *AppInfo {
	once.Do(func() {
		instance = &AppInfo{
			ServiceBaseURL: "/_service/v0",
			ServiceVersion: "/v0",
			ApiBaseURL:     "/v0",
			ApiVersion:     "v0",
			ServiceLocator: "//banckend.core",
			ReleaseInfo: &ReleaseInfo{
				Date:    "2024-10-02T18:09:51.231Z",
				Commit:  "uuid",
				Name:    "backendcore",
				Uuid:    "commit",
				Version: "0.34.1-SNAPSHOT",
			},
		}
	})
	return instance
}
