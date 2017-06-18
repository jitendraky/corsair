package config

import (
  "io/ioutil"
  "os"
  "log"

  "corsair/corsair"

  "framework/version"

  "gopkg.in/yaml.v2"
)

type State struct {
  DebugMode          bool             `yaml:"debug"`
  Quiet              bool             `yaml:"quiet"`
  MaxCPU             string           `yaml:"max_cpu"`
  LogPath            string           `yaml:"log_path"`
  PIDPath            string           `ymal:"pid_path"`
  Server            struct{
    ServerType      string             `yaml:"server_type"`

    UserAgent       string             `yaml:"user_agent"`
  } `yaml:"server"`
  Databases         struct{
    MemoryDatabase  struct{
      Enabled         bool             `yaml"enabled"`
    } `yaml:"memory_database"`
    FileDatabase    struct{
      Enabled         bool             `yaml:"enabled"`
      DatabasePath    string           `yaml:"database_path"`
    } `yaml:"file_database"`
  } `yaml:"databases"`
  Plugins           []Plugin           `yaml:"plugins"`
}

type Plugin struct {
  Name                string
  ImportPath          string
  Repository          string
  LockCommit          string
  Version             version.Version
  Enabled             bool
  Directive
}

type Directive struct {
  Name                string           `yaml:"name"`
  Weight              string           `yaml:"weight"`
}

func Load(filePath string) (state State) {
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

// TODO: Possibly combine this with the above YAML, for a single unified config
// by a consistent and not proprietary config language.
func CorsairFileLoader(configPath, serverType string) (corsair.Input) {
  var contents []byte
  var err error
  if configPath == "" {
    contents, err = ioutil.ReadFile(corsair.DefaultConfigFile)
  } else if configPath == "stdin" {
  	fi, err := os.Stdin.Stat()
  	if err == nil && fi.Mode()&os.ModeCharDevice == 0 {
  		contents, err = ioutil.ReadAll(os.Stdin)
      configPath = os.Stdin.Name()
  	}
  } else {
    contents, err = ioutil.ReadFile(configPath)
  }
  if err != nil {
    return nil
  }
  return corsair.CorsairfileInput{
    Contents:       []byte(contents),
    Filepath:       configPath,
    ServerTypeName: serverType,
  }
}

