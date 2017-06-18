package config

import (
  "io/ioutil"
  "os"
  "log"

  "./version"

  "gopkg.in/yaml.v2"
)

type WatchHook struct {
  Name          string                 `yaml:"name"`
  Paths         []string               `yaml:"paths"`
  Recursive     bool                   `yaml:"recursive"`
  Extensions    string                 `yaml:"watch"`
  Command       string                 `yaml:"command"`
}

type Language struct {
  Name          string                 `yaml:"name"`
  Path          string                 `yaml:"path"`
  Console       string                 `yaml:"console"`
  Version       version.Version        `yaml:"version"`
}

type ProjectComponent struct {
  Name          string                 `yaml:"name"`
  Version       version.Version        `yaml:"version"`
  Language                             `yaml:"language"`
  Repository    SourceRepository       `yaml:"repository"`
}

type SourceRepository struct {
  Repository    string                  `yaml:"repository"`
  Type          string                  `yaml:"type"`
}

type Command struct {
  Name              string             `yaml:"name"`
  WorkingPath       string             `yaml:"working_path"`
  Command           string             `yaml:"command"`
}

type State struct {
  DebugMode         bool               `yaml:"debug"`
  AfterStartup      string             `yaml:"after_startup"`
  BeforeShutdown    string             `yaml:"before_shutdown"`
  WorkingPath       string             `yaml:"working_path"`
  Project         struct {
    Name            string             `yaml:"name"`
    Version         version.Version    `yaml:"version"`
    Environment     string             `yaml:"environment"`
    Language
    Repository      string             `yaml:"repository"`
    Path            string             `yaml:"path"`
    Components      []ProjectComponent `yaml:"components"`
  } `yaml:"project"`
  Commands          []Command          `yaml:"commands"`
  Database        struct {
    DatabasePath      string           `yaml:"database_path"`
    SearchIndexPath   string           `yaml:"search_index_path"`
  } `yaml:"database"`
  Watch           struct {
    Enabled         bool               `yaml:"enabled"`
    Hooks           []WatchHook        `yaml:"hooks"`
  } `yaml:"watch"`
  WebUI           struct {
    Enabled         bool               `yaml:"enabled"`
    Authentication  bool               `yaml:"authentication"`
    APIVersion      string             `yaml:"api_version"`
    Host            string             `yaml:"host"`
    Port            string             `yaml:"port"`
  } `yaml:"web_ui"`
}

func Load(filePath string) (state State) {
  state.WorkingPath, _ = os.Executable()
  //TODO: Check if config exists, if it does not, write a default configuration
  //state = State{
  //  Watch: Watch{
  //  }
  //}
  configFile, err := ioutil.ReadFile(filePath)
  err = yaml.Unmarshal(configFile, &state)
  if err != nil {
    log.Println("[Error]: ", err)
    // init the config
    // then recall Load
  }

  return state
}
