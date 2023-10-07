package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/samber/lo"
	"strings"
	"sync"
	"time"
	"wox/util"
)

type StoreManifest struct {
	Name string
	Url  string
}

type StorePluginManifest struct {
	Id             string
	Name           string
	Author         string
	Version        string
	Runtime        Runtime
	Description    string
	IconUrl        string
	Website        string
	DownloadUrl    string
	ScreenshotUrls []string
	DateCreated    string
	DateUpdated    string
}

var storeInstance *Store
var storeOnce sync.Once

type Store struct {
	pluginManifests []StorePluginManifest
}

func GetStoreManager() *Store {
	storeOnce.Do(func() {
		storeInstance = &Store{}
	})
	return storeInstance
}

func (s *Store) GetStoreManifests(ctx context.Context) []StoreManifest {
	return []StoreManifest{
		{
			Name: "Wox Official Plugin",
			Url:  "https://raw.githubusercontent.com/Wox-launcher/Wox/v2/plugin-store.json",
		},
	}
}

// get plugin manifests from plugin stores, and update in the background every 10 minutes
func (s *Store) Start(ctx context.Context) {
	s.pluginManifests = s.GetStorePluginManifests(ctx)

	util.Go(ctx, "load store plugins", func() {
		for range time.NewTicker(time.Minute * 10).C {
			pluginManifests := s.GetStorePluginManifests(util.NewTraceContext())
			if len(pluginManifests) > 0 {
				s.pluginManifests = pluginManifests
			}
		}
	})
}

func (s *Store) GetStorePluginManifests(ctx context.Context) []StorePluginManifest {
	var storePluginManifests []StorePluginManifest

	for _, store := range s.GetStoreManifests(ctx) {
		storePluginManifest, manifestErr := s.GetStorePluginManifest(ctx, store)
		if manifestErr != nil {
			logger.Error(ctx, fmt.Sprintf("failed to get plugin manifest from %s store: %s", store.Name, manifestErr.Error()))
			continue
		}

		for _, manifest := range storePluginManifest {
			existingManifest, found := lo.Find(storePluginManifests, func(manifest StorePluginManifest) bool {
				return manifest.Id == manifest.Id
			})
			if found {
				existingVersion, existingErr := semver.NewVersion(existingManifest.Version)
				currentVersion, currentErr := semver.NewVersion(manifest.Version)
				if existingErr != nil && currentErr != nil {
					if existingVersion.GreaterThan(currentVersion) {
						logger.Info(ctx, fmt.Sprintf("skip %s(%s) from %s store, because it's already installed(%s)", manifest.Name, manifest.Version, store.Name, existingManifest.Version))
						continue
					}
				}
			}

			storePluginManifests = append(storePluginManifests, manifest)
		}
	}

	logger.Info(ctx, fmt.Sprintf("found %d plugins from stores", len(storePluginManifests)))
	return storePluginManifests
}

func (s *Store) GetStorePluginManifest(ctx context.Context, store StoreManifest) ([]StorePluginManifest, error) {
	logger.Info(ctx, fmt.Sprintf("start to get plugin manifest from %s(%s)", store.Name, store.Url))

	response, getErr := util.Get(store.Url)
	if getErr != nil {
		return nil, getErr
	}

	var storePluginManifests []StorePluginManifest
	unmarshalErr := json.Unmarshal(response, &storePluginManifests)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return storePluginManifests, nil
}

func (s *Store) Search(ctx context.Context, keyword string) []StorePluginManifest {
	return lo.Filter(s.pluginManifests, func(manifest StorePluginManifest, _ int) bool {
		return strings.Contains(strings.ToLower(manifest.Name), strings.ToLower(keyword))
	})
}