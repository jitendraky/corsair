package main

import (
	//"errors"
	"fmt"
  //"time"
	//"io/ioutil"
	"log"
	"os"
  "os/exec"
	//"runtime"
	//"strconv"
	//"strings"

	//"gopkg.in/natefinch/lumberjack.v2"
	"github.com/xenolf/lego/acme"
  cli "gopkg.in/alecthomas/kingpin.v2"

  // ## Corsair Framework ##
  "framework/version"
  "framework/config"
  "framework/database/memory"

  // ## *Library Defined Plugins* ##
	"corsair/corsair"
	_ "corsair/corsairhttp"
	"corsair/corsairtls"

  // ## *User Defined Plugins* ##
  // ## Include plugins here, with the special
  // ## character '_' ifront of the imported library.
  _ "plugins/boilerplate"
  _ "plugins/corsair-search"
)

type Application struct {
  Name              string
  Input             corsair.Input
  ConfigPath        string
  Config            config.State
  MemoryDB          *memory.DB
  Process           *exec.Cmd
  MaxCPU            string
  PIDFilePath       string
  LogFilePath       string
  Event             chan string
  Flags struct{
    Quiet           string
	  ConfigFile      string
	  MaxCPU          string
	  LogPath         string
    PIDPath         string
  }
}

func (self *Application) flagsAndCommands(){
  args           := cli.New("corsair", "Corsair application framework")
  Quiet          :=  args.Flag("quiet", "List installed plugins, not necessarily actived plugins.").Bool()
  ConfigFile     :=  args.Flag("config", "Corsairfile to load, default location is the current working directory.").String()
  MaxCPU         :=  args.Flag("cpu", "Maximum CPU allowed for the Corsair application process.").Default("100%").String()
  LogPath        :=  args.Flag("log", "Path for logfile, default is std.out.").String()
  PIDPath        :=  args.Flag("pid", "Path for PIDfile.").String()

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

  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
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

func (self *Application)printBanner() {
  fmt.Println("  _____                    ")
  fmt.Println(" /     \\                   ")
  fmt.Println("| () () __________________ ")
  fmt.Println(" \\  ^  /                  |")
  fmt.Println("  | |", self.Name, "/",)
  fmt.Println("  |||||___________________|")
  fmt.Println("              ", self.Build.Version.ToString(), "")
}

func main() {
  //
	corsair.TrapSignals()
	corsair.RegisterCorsairfileLoader("flag", corsair.LoaderFunc(confLoader))
	corsair.SetDefaultCorsairfileLoader("default", corsair.LoaderFunc(defaultLoader))
  //
  app := Application{
    Name:             "Corsair",
    WorkingDirectory: os.Executable(),
    Config:           config.Load(os.Executable()+"corsair.yaml"),
    MemoryDB:         memory.Open(":memory:"),
    TLS: TLS{
      CAUrl: "https://acme-v01.api.letsencrypt.org/directory",
      Email: "Default ACME CA account email address",
    },
    Build: Build{
      Version: version.Version{Major: 0, Minor: 1, Patch: 0},
      Development: true,
    },
  }
  app.printBanner()
	app.flagsAndCommands()

  // TODO: This does not seem right at all.
	acme.UserAgent = app.Config.Server.UserAgent

	// Executes Startup events
	corsair.EmitEvent(corsair.StartupEvent, nil)

	// Get Corsairfile inputA
  if app.Flags.ConfigPath != "" {
	  corsairfileinput = corsair.LoadCorsairfile(app.Config.serverType)
  } else if app.Config.ConfigPath {
	  corsairfileinput = corsair.LoadCorsairfile(app.Config.serverType)
  }
	// Start your engines
	instance, err := corsair.Start(corsairfileinput)
	if err != nil {
		mustLogFatalf("%v", err)
	}

	// Twiddle your thumbs
	instance.Wait()
}

