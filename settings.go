package main // import "git.zeuz.io/devel/olympuz/sdk/golang/zeuztool/cli"

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"git.zeuz.io/sdk/golang/zeuzsdk"
)

// SettingsCtx doctodo
type SettingsCtx struct {
	CtxName      string
	APIEndpoint  string
	DevID        zeuzsdk.DeveloperID
	DevLogin     string
	DevPassword  string
	ProjID       zeuzsdk.ProjID
	EnvID        zeuzsdk.EnvID
	AccountID    zeuzsdk.AccountID
	TeamID       zeuzsdk.TeamID
	UserID       zeuzsdk.UserID
	SessionID    zeuzsdk.SessionID
	SessionKey   string
	EMail        string
	APIKey       string
	Title        string
	PersistLogin bool
}

// SettingsList doctodo
type SettingsList struct {
	CtxList []*SettingsCtx
	Default string
}

// Settings doctodo
var Settings SettingsList

// CurSettings doctodo
var CurSettings *SettingsCtx
var settingsStoreCS uint64
var settingsFilename string

func GetSettingsFilename() string {
	if settingsFilename != "" {
		return settingsFilename
	}
	var localConfig string

	switch runtime.GOOS {
	case "windows":
		localConfig = os.Getenv("APPDATA")
	case "darwin":
		localConfig = os.Getenv("HOME") + "/Library/Application Support"
	default:
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			localConfig = os.Getenv("XDG_CONFIG_HOME")
		} else {
			localConfig = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}
	return filepath.Join(localConfig, "zeuz.io/zeuzcmd.json")
}

func settingsGetStoreCS() uint64 {
	data, _ := json.MarshalIndent(&Settings, "", "\t")
	xHash := fnv.New64()
	xHash.Write(data)
	return xHash.Sum64()
}

func SettingsClear() {
	CurSettings = nil
	Settings.CtxList = nil
	Settings.Default = "default"
}

// SettingsLoad doctodo
func SettingsLoad() error {
	fmt.Printf("Load LocalConfig: %v\n", GetSettingsFilename())
	data, err := ioutil.ReadFile(GetSettingsFilename())

	if err != nil {
		if os.IsNotExist(err) {
			SettingsSelect("default")
			SettingsSetDefault("default")
			SettingsSave()
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, &Settings)
	if err == nil {
		SettingsSelect(Settings.Default)
	}
	settingsStoreCS = settingsGetStoreCS()
	return err
}

func GetCurSettings() *SettingsCtx {
	if CurSettings == nil {
		SettingsLoad()
	}
	return CurSettings
}

// SettingsSave doctodo
func SettingsSave() error {
	if settingsStoreCS == settingsGetStoreCS() {
		return nil
	}
	data, err := json.MarshalIndent(&Settings, "", "\t")
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(GetSettingsFilename()), os.ModePerm)
	//fmt.Printf("CreateOrUpdate LocalConfig: %v\n", settingsfilename)
	err = ioutil.WriteFile(GetSettingsFilename(), data, os.ModePerm)
	return err
}

// SettingsSetDefault doctodo
func SettingsSetDefault(name string) {
	SettingsSelect(name)
	Settings.Default = name
}

// SettingsSelect doctodo
func SettingsSelect(name string) {
	for _, c := range Settings.CtxList {
		if c.CtxName == name {
			CurSettings = c
			return
		}
	}
	c := &SettingsCtx{
		CtxName: name,
	}
	Settings.CtxList = append(Settings.CtxList, c)
	CurSettings = c
}

// SettingsRemove doctodo
func SettingsRemove(name string) {
	if name == Settings.Default {
		SettingsSetDefault("default")
	}
	for i, c := range Settings.CtxList {
		if c.CtxName == name {
			Settings.CtxList = append(Settings.CtxList[:i], Settings.CtxList[i+1:]...)
			return
		}
	}
}
