package main

import (
  "fmt"
  "log"
  "strings"
  "time"
  "os"
  "os/signal"
  "os/exec"
  "sync"
  "syscall"

  // TODO: If these were put into the src folder,
  // then the ./ deleclaration would be avoided.
  //"./config/shutdown"
  //"./models/file"
  //"./watch"
  //"./database"
  "./config"
  "./config/version"
  "./watcher"

  //shell "github.com/abiosoft/ishell"
  "gopkg.in/alecthomas/kingpin.v2"
)

type Flags struct {
    ProjectPath     *string
    Watching        *bool
    WatchExtensions *string
    WatchPath       *string
    WatchRecursive  *bool
    WatchCommand    *string
    WebUI           *bool
    WebUIPort       *string
    WebUIHost       *string
}

type Application struct {
  Name            string
  Version         version.Version
  Config          config.State
  RunningProcess  *exec.Cmd
  ExitChannel     chan bool
  WaitGroup       sync.WaitGroup
  //Database      database.Context
}

func cacheAllFiles(fileDataSet [](*FileData)) error {
  //TODO: Stick into the File BoltDB here
  for _, fileData := range fileDataSet {
    _, _ = fileData.GetContents(true, false)
  }
  return nil
}

func (self Application) printBanner() {
  fmt.Println("Gravity Build System v", self.Version.ToString() )
}

func (self *Application) WatchPath(path string, command string, recursive bool, milliseconds int){
  //TODO: Move to initWatch in watcher.go, then just pass in the functions you want to hook
  // one vents
  w := watcher.New()
	w.SetMaxEvents(1)
  // TODO: Pass the events to watch. Only notify rename and move events.
	//w.FilterOps(watcher.Rename, watcher.Move)
	go func() {
		for {
			select {
			case event := <-w.Event:
        // Print the event's info.
        if self.RunningProcess != nil {
          KillProcess(self.RunningProcess)
        }
        fmt.Println("Attempting to run command, ", command)
        self.RunningProcess = StartCommand(path, command)
				fmt.Println(event)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()
  if recursive {
    fmt.Println("Adding path recursively, this takes time, please wait...")
	  if err := w.AddRecursive(path); err != nil {
	    log.Fatalln(err)
    }
  }else{
    if err := w.Add(path); err != nil {
	    log.Fatalln(err)
	  }
  }
	//for path, f := range w.WatchedFiles() {
	//  fmt.Printf("%s: %s\n", path, f.Name())
  //}
	if err := w.Start(time.Millisecond * time.Duration(milliseconds)); err != nil {
		log.Fatalln(err)
	}
}

func (self *Application) signalHandler(){
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
      for {
        s := <-c
        fmt.Println("Got signal:", s)
        if self.RunningProcess != nil {
          KillProcess(self.RunningProcess)
        }
        self.ExitChannel <- true
        os.Exit(0)
      }
  }()
}

func (self *Application) processCommands() {
  //switch kingpin.MustParse(app.Parse(os.Args[1:])) {
  //case ""
  //}
  for _, command := range self.Config.Commands {
		fmt.Println("Command would be registered: ", command.Name)
    //post        = kingpin.Command("post", "Post a message to a channel.")
  }
}

// TODO: Migrate from a flag based to a command based system
// https://github.com/urfave/cli may be a better choice over kingpin
func (self *Application) processFlags() {
  flags := Flags{
      ProjectPath:      kingpin.Flag("path", "The directory to serve").String(),
      Watching:         kingpin.Flag("watch", "Enable file path watching, triggering an event on change").Bool(),
      WatchPath:        kingpin.Flag("wpath", "File extensions to watch for change in the project directory").String(),
      WatchExtensions:  kingpin.Flag("wext", "File extensions to watch for change in the project directory").String(),
      WatchCommand:     kingpin.Flag("wcmd", "Command to trigger on change of watched files").String(),
      WebUI:            kingpin.Flag("webui", "Enable web UI").Bool(),
      WebUIHost:        kingpin.Flag("host", "Host/interface to serve the WebUI").String(),
      WebUIPort:        kingpin.Flag("port", "Port to serve the WebUI").String(),
  }
  kingpin.Parse()
  if *flags.ProjectPath != "" {
    self.Config.Project.Path = *flags.ProjectPath
  }
  if *flags.Watching != false {
    self.Config.Watch.Enabled = *flags.Watching
    if *flags.WatchPath != "" {
      // TODO: Count the number of commas in the WatchPath add multiple paths
      paths := strings.Split(*flags.WatchPath, ",")
      self.Config.Watch.Hooks = append(self.Config.Watch.Hooks, config.WatchHook{Paths: paths, Extensions: *flags.WatchExtensions, Command: *flags.WatchCommand})
    }
  }
  if *flags.WebUI != false {
    self.Config.WebUI.Enabled = *flags.WebUI
    if *flags.WebUIHost != "" {
      self.Config.WebUI.Host = *flags.WebUIHost
    }
    if *flags.WebUIPort != "" {
      self.Config.WebUI.Port = *flags.WebUIPort
    }
  }else{
    if self.Config.WebUI.Host == "" {
      self.Config.WebUI.Host = "127.0.0.1"
    }
    if self.Config.WebUI.Port == "" {
      self.Config.WebUI.Port = "8005"
    }
  }
  if *flags.WebUI {
    self.Config.WebUI.Enabled = true
    self.Config.WebUI.Host    = *flags.WebUIHost
    self.Config.WebUI.Port    = *flags.WebUIPort
  }
}

func main() {
  // Initialization & Configuration
  // TODO: If the configuration doesn't exist, populate a new one with reasonable defaults
  app := Application{
    Name:     "Gravity Build System",
    Version:  version.Version{Major: 0, Minor: 1, Patch: 0},
    Config:   config.Load("./build.yaml"),
    ExitChannel: make(chan bool, 1),
    //Database: database.Context{
    //  Path:   "~/.local/gravity/database.db",
    //  Opened: false,
    //},
  }
  app.processFlags()
  app.processCommands()
  app.printBanner()
  // Signal Handling
  app.signalHandler()
  // Database
  //app.Database.Open(app.Database.Path, true)
  //defer app.Database.Close()
  //app.Database.Indexes["files"] = app.Database.InitiateSearchIndex("~/.local/gravity/files.bleve")
  // Set GOPATH
  if app.Config.AfterStartup != "" {
    fmt.Println("Executing build tool startup up command...")
    fmt.Println(app.Config.AfterStartup)
    StartCommand(app.Config.Project.Path, app.Config.AfterStartup)
  }
  app.WaitGroup.Add(1)
  // Watcher
  if app.Config.Watch.Enabled {
    fmt.Println("[Watch Enabled]")
    if len(app.Config.Watch.Hooks) > 0 {
      for _, hook := range app.Config.Watch.Hooks {
        fmt.Println("  Watch recursive: ", hook.Recursive)
        fmt.Println("  Watch command  : ", hook.Command)
        for _, path := range hook.Paths {
          fmt.Println("  Watch Path     : ", path)
          app.WatchPath(path, hook.Command, hook.Recursive, 100)
        }
      }
    }
  }
  // WebUI
  fmt.Println("WebUI Enabled? ", app.Config.WebUI.Enabled)
  if app.Config.WebUI.Enabled {
    go app.InitiateWebUI()
  }
  // Console UI
  // TODO: This one may be too opinionated as it gets mixed up with our signalsm in this file
  //shell := shell.New()
  //shell.Println("Gravity Build System Shell")
  //shell.Run()

  

  fmt.Println("awaiting signal")
  <-app.ExitChannel
}
