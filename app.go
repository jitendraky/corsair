package main

import (
	"fmt"
	"log"
	"os"
	//"gopkg.in/natefinch/lumberjack.v2"
	"github.com/xenolf/lego/acme"
	cli "gopkg.in/alecthomas/kingpin.v2"
	// ## Corsair Framework ##
	"corsair/config"
	"corsair/config/version"
	"corsair/models"
	"corsair/database/memory"
	"corsair/corsair"
	_ "corsair/network/http"
	"corsair/network/tls"
	// ## *User Defined Plugins* ##
	// ## Include plugins here, with the special
	// ## character '_' ifront of the imported library.
	_ "plugins/boilerplate"
	_ "plugins/corsair-search"
)

type Daemon struct {
  // Using embedding principle to inherit the corsair application object
  App       corsair.Application
}

func (self *Process) flagsAndCommands(){
  args           := cli.New("corsair", "Corsair application framework")
  quiet          :=  args.Flag("quiet", "List installed plugins, not necessarily actived plugins.").Bool()
  configPath     :=  args.Flag("config", "Corsairfile to load, default location is the current working directory.").String()
  maxCPU         :=  args.Flag("cpu", "Maximum CPU allowed for the Corsair application process.").Default("100%").String()
  logPath        :=  args.Flag("log", "Path for logfile, default is std.out.").String()
  pidPath        :=  args.Flag("pid", "Path for PIDfile.").String()
  // TODO: Commands should be in a app-cli executable, and the daemon should be in a appd executable
  plugins        :=  args.Command("plugins", "List installed plugins, not necessarily actived plugins.")
  pluginsAction  :=  plugins.Arg("action", "Available actions on all plugins include: list, update, install.")
  plugin         :=  args.Command("plugin", "List installed plugins, not necessarily actived plugins.")
  pluginAction   :=  plugin.Arg("action", "Available actions include: list, update, install.")
  pluginItem     :=  plugin.Arg("item", "Available actions include: list, update, install.")
  revoke         :=  args.Command("revoke", "Revoke a certificate associated with specified domain.")
  revokeHost     :=  revoke.Arg("host", "All certiciates with this hostname will be removed.")
  validate       := args.Command("validate", "Parse the Corsairfile without starting the server")
  version        :=  args.Command("version", "Show version and other build information of the application.")
  versionFlag    :=  args.Flag("version", "Show version and other build information of the application.").Bool()
  if *versionFlag {
    self.printBanner()
    os.Exit(0)
  }
  if *configPath != "" {
    self.App.ConfigPath = *configPath
  }
  if *logPath != "" {
    self.App.Config.LogPath = *logPath
  }
  if *pidPath != "" {
    self.App.Config.PIDPath = *pidPath
  }
  if *maxCPU != "" {
    self.App.Config.MaxCPU = *maxCPU
  }
  if *quiet {
    self.App.Config.Quiet = *quiet
  }

  switch cli.MustParse(args.Parse(os.Args[1:])) {
  case plugins.FullCommand():
		fmt.Println(corsair.DescribePlugins())
		os.Exit(0)
  case plugin.FullCommand():
    fmt.Println("[Error] Currently unimplemented plugin")
  case revoke.FullCommand():
    fmt.Println("CAN YOU USE ", revokeHost, " instead of os.Args[2]?")
		err := corsairtls.Revoke(os.Args[2])
		if err != nil {
			mustLogFatalf("%v", err)
		}
		fmt.Printf("Revoked certificate for %s\n", os.Args[2])
		os.Exit(0)
  case version.FullCommand():
    self.printBanner()
    os.Exit(0)
  case validate.FullCommand():
		err := corsair.ValidateAndExecuteDirectives(self.Input, nil, true)
		if err != nil {
			mustLogFatalf("%v", err)
		}
		msg := "Corsairfile is valid"
		fmt.Println(msg)
		log.Printf("[INFO] %s", msg)
		os.Exit(0)
  }
}

func (self *Daemon)printBanner() {
  fmt.Println("  _____                    ")
  fmt.Println(" /     \\                   ")
  fmt.Println("| () () __________________ ")
  fmt.Println(" \\  ^  /                  |")
  fmt.Println("  | |", self.Name, "/",)
  fmt.Println("  |||||___________________|")
  fmt.Println("              ", self.Build.Version.ToString(), "")
}

func main() {
	//corsair.RegisterCorsairfileLoader("flag", corsair.LoaderFunc(confLoader))
	//corsair.SetDefaultCorsairfileLoader("default", corsair.LoaderFunc(defaultLoader))
  workingPath, _ := os.Executable()
  daemon = Daemon{
    App:   corsair.Application{
      Name:             "Corsair",
      WorkingPath:      workingPath,
      Config:           config.Load(workingPath+"corsair.yaml"),
      MemoryDB:         memory.Open(":memory:"),
      TLS:  models.TLS{
        CAUrl: "https://acme-v01.api.letsencrypt.org/directory",
        Email: "Default ACME CA account email address",
      },
      Build: models.Build{
        Version: version.Version{Major: 0, Minor: 1, Patch: 0},
	    Development: true,
      },
    },
  }
	daemon.App.TrapSignals()
  daemon.printBanner()
	daemon.flagsAndCommands()
	acme.UserAgent = daemon.App.Config.Server.UserAgent

	corsair.EmitEvent(corsair.StartupEvent, nil)
	daemon.App.Input = config.LoadCorsairfile(daemon.App.Config.Server.ServerType)

	instance, err := corsair.Start(corsairfileinput)
	if err != nil {
		mustLogFatalf("%v", err)
	}

	instance.Wait()
}

