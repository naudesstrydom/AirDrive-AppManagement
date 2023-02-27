package service

import (
	"fmt"
	"strings"

	"github.com/IceWhaleTech/CasaOS-AppManagement/codegen"
	"github.com/IceWhaleTech/CasaOS-AppManagement/pkg/config"
	"github.com/IceWhaleTech/CasaOS-Common/utils"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AppStoreManagement struct {
	onAppStoreRegister   []func(string) error
	onAppStoreUnregister []func(string) error

	appStoreMap     map[string]*AppStore
	defaultAppStore *AppStore
}

func (a *AppStoreManagement) AppStoreList() []codegen.AppStoreMetadata {
	return lo.Map(config.ServerInfo.AppStoreList, func(appStoreURL string, id int) codegen.AppStoreMetadata {
		return codegen.AppStoreMetadata{
			ID:  &id,
			URL: &appStoreURL,
		}
	})
}

func (a *AppStoreManagement) OnAppStoreRegister(fn func(string) error) {
	a.onAppStoreRegister = append(a.onAppStoreRegister, fn)
}

func (a *AppStoreManagement) OnAppStoreUnregister(fn func(string) error) {
	a.onAppStoreUnregister = append(a.onAppStoreUnregister, fn)
}

func (a *AppStoreManagement) RegisterAppStore(appstoreURL string) (*codegen.AppStoreMetadata, error) {
	appstoreURL = strings.ToLower(appstoreURL)

	// check if appstore already exists
	config.ServerInfo.AppStoreList = lo.Map(config.ServerInfo.AppStoreList,
		func(url string, id int) string {
			return strings.ToLower(url)
		})

	for i, url := range config.ServerInfo.AppStoreList {
		if url == appstoreURL {
			return &codegen.AppStoreMetadata{
				ID:  &i,
				URL: &config.ServerInfo.AppStoreList[i],
			}, nil
		}
	}

	appstore, err := NewAppStore(appstoreURL)
	if err != nil {
		return nil, err
	}

	if err := appstore.UpdateCatalog(); err != nil {
		return nil, err
	}

	// if everything is good, add to the list
	a.appStoreMap[appstoreURL] = appstore

	config.ServerInfo.AppStoreList = append(config.ServerInfo.AppStoreList, appstoreURL)

	if err := config.SaveSetup(); err != nil {
		return nil, err
	}

	for _, fn := range a.onAppStoreRegister {
		if err := fn(appstoreURL); err != nil {
			return nil, err
		}
	}

	return &codegen.AppStoreMetadata{
		ID:  utils.Ptr(len(config.ServerInfo.AppStoreList) - 1),
		URL: &appstoreURL,
	}, nil
}

func (a *AppStoreManagement) UnregisterAppStore(appStoreID uint) error {
	appStoreURL := config.ServerInfo.AppStoreList[appStoreID]

	config.ServerInfo.AppStoreList = append(config.ServerInfo.AppStoreList[:appStoreID], config.ServerInfo.AppStoreList[appStoreID+1:]...)

	if err := config.SaveSetup(); err != nil {
		return err
	}

	delete(a.appStoreMap, appStoreURL)

	for _, fn := range a.onAppStoreUnregister {
		if err := fn(appStoreURL); err != nil {
			return err
		}
	}
	return nil
}

func (a *AppStoreManagement) Catalog() map[string]*ComposeApp {
	catalog := map[string]*ComposeApp{}

	for _, appStore := range a.appStoreMap {
		for appStoreID, composeApp := range appStore.Catalog() {
			catalog[appStoreID] = composeApp
		}
	}

	if len(catalog) == 0 {
		logger.Info("No appstore registered")
		if a.defaultAppStore == nil {
			logger.Info("WARNING - no default appstore")
			return map[string]*ComposeApp{}
		}

		logger.Info("Using default appstore")
		return a.defaultAppStore.Catalog()
	}

	return catalog
}

func (a *AppStoreManagement) UpdateCatalog() {
	for url, appStore := range a.appStoreMap {
		if err := appStore.UpdateCatalog(); err != nil {
			logger.Error("error while updating catalog for app store", zap.Error(err), zap.String("url", url))
		}
	}
}

func (a *AppStoreManagement) ComposeApp(id string) *ComposeApp {
	for _, appStore := range a.appStoreMap {
		if composeApp := appStore.ComposeApp(id); composeApp != nil {
			return composeApp
		}
	}

	logger.Info("No appstore registered")

	if a.defaultAppStore == nil {
		logger.Info("WARNING - no default appstore")
		return nil
	}

	logger.Info("Using default appstore")

	return a.defaultAppStore.ComposeApp(id)
}

func NewAppStoreManagement() *AppStoreManagement {
	defaultAppStore, err := NewDefaultAppStore()
	if err != nil {
		fmt.Printf("error while loading default appstore: %s\n", err.Error())
	}

	appStoreManagement := &AppStoreManagement{
		appStoreMap:     map[string]*AppStore{},
		defaultAppStore: defaultAppStore,
	}

	return appStoreManagement
}
