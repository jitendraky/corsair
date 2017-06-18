package corsair

import (
	"fmt"
	"log"
	"net"
	"sort"

	"corsair/corsairfile"
)

// These are all the registered plugins.
var (
	// serverTypes is a map of registered server types.
	serverTypes = make(map[string]ServerType)

	// plugins is a map of server type to map of plugin name to
	// Plugin. These are the "general" plugins that may or may
	// not be associated with a specific server type. If it's
	// applicable to multiple server types or the server type is
	// irrelevant, the key is empty string (""). But all plugins
	// must have a name.
	plugins = make(map[string]map[string]Plugin)

	// eventHooks is a map of hook name to Hook. All hooks plugins
	// must have a name.
	eventHooks = make(map[string]EventHook)

	// parsingCallbacks maps server type to map of directive
	// to list of callback functions. These aren't really
	// plugins on their own, but are often registered from
	// plugins.
	parsingCallbacks = make(map[string]map[string][]ParsingCallback)

	// corsairfileLoaders is the list of all Corsairfile loaders
	// in registration order.
	corsairfileLoaders []corsairfileLoader
)

// DescribePlugins returns a string describing the registered plugins.
func DescribePlugins() string {
	str := "Server types:\n"
	for name := range serverTypes {
		str += "  " + name + "\n"
	}

	// List the loaders in registration order
	str += "\nCorsairfile loaders:\n"
	for _, loader := range corsairfileLoaders {
		str += "  " + loader.name + "\n"
	}
	if defaultCorsairfileLoader.name != "" {
		str += "  " + defaultCorsairfileLoader.name + "\n"
	}

	if len(eventHooks) > 0 {
		// List the event hook plugins
		str += "\nEvent hook plugins:\n"
		for hookPlugin := range eventHooks {
			str += "  hook." + hookPlugin + "\n"
		}
	}

	// Let's alphabetize the rest of these...
	var others []string
	for stype, stypePlugins := range plugins {
		for name := range stypePlugins {
			var s string
			if stype != "" {
				s = stype + "."
			}
			s += name
			others = append(others, s)
		}
	}

	sort.Strings(others)
	str += "\nOther plugins:\n"
	for _, name := range others {
		str += "  " + name + "\n"
	}

	return str
}

// ValidDirectives returns the list of all directives that are
// recognized for the server type serverType. However, not all
// directives may be installed. This makes it possible to give
// more helpful error messages, like "did you mean ..." or
// "maybe you need to plug in ...".
func ValidDirectives(serverType string) []string {
	stype, err := getServerType(serverType)
	if err != nil {
		return nil
	}
	return stype.Directives()
}

// ServerListener pairs a server to its listener and/or packetconn.
type ServerListener struct {
	server   Server
	listener net.Listener
	packet   net.PacketConn
}

// LocalAddr returns the local network address of the packetconn. It returns
// nil when it is not set.
func (s ServerListener) LocalAddr() net.Addr {
	if s.packet == nil {
		return nil
	}
	return s.packet.LocalAddr()
}

// Addr returns the listener's network address. It returns nil when it is
// not set.
func (s ServerListener) Addr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Context is a type which carries a server type through
// the load and setup phase; it maintains the state
// between loading the Corsairfile, then executing its
// directives, then making the servers for Corsair to
// manage. Typically, such state involves configuration
// structs, etc.
type Context interface {
	// Called after the Corsairfile is parsed into server
	// blocks but before the directives are executed,
	// this method gives you an opportunity to inspect
	// the server blocks and prepare for the execution
	// of directives. Return the server blocks (which
	// you may modify, if desired) and an error, if any.
	// The first argument is the name or path to the
	// configuration file (Corsairfile).
	//
	// This function can be a no-op and simply return its
	// input if there is nothing to do here.
	InspectServerBlocks(string, []corsairfile.ServerBlock) ([]corsairfile.ServerBlock, error)

	// This is what Corsair calls to make server instances.
	// By this time, all directives have been executed and,
	// presumably, the context has enough state to produce
	// server instances for Corsair to start.
	MakeServers() ([]Server, error)
}

// RegisterServerType registers a server type srv by its
// name, typeName.
func RegisterServerType(typeName string, srv ServerType) {
	if _, ok := serverTypes[typeName]; ok {
		panic("server type already registered")
	}
	serverTypes[typeName] = srv
}

// ServerType contains information about a server type.
type ServerType struct {
	// Function that returns the list of directives, in
	// execution order, that are valid for this server
	// type. Directives should be one word if possible
	// and lower-cased.
	Directives func() []string

	// DefaultInput returns a default config input if none
	// is otherwise loaded. This is optional, but highly
	// recommended, otherwise a blank Corsairfile will be
	// used.
	DefaultInput func() Input

	// The function that produces a new server type context.
	// This will be called when a new Corsairfile is being
	// loaded, parsed, and executed independently of any
	// startup phases before this one. It's a way to keep
	// each set of server instances separate and to reduce
	// the amount of global state you need.
	NewContext func() Context
}

// Plugin is a type which holds information about a plugin.
type Plugin struct {
	// ServerType is the type of server this plugin is for.
	// Can be empty if not applicable, or if the plugin
	// can associate with any server type.
	ServerType string

	// Action is the plugin's setup function, if associated
	// with a directive in the Corsairfile.
	Action SetupFunc
}

// RegisterPlugin plugs in plugin. All plugins should register
// themselves, even if they do not perform an action associated
// with a directive. It is important for the process to know
// which plugins are available.
//
// The plugin MUST have a name: lower case and one word.
// If this plugin has an action, it must be the name of
// the directive that invokes it. A name is always required
// and must be unique for the server type.
func RegisterPlugin(name string, plugin Plugin) {
	if name == "" {
		panic("plugin must have a name")
	}
	if _, ok := plugins[plugin.ServerType]; !ok {
		plugins[plugin.ServerType] = make(map[string]Plugin)
	}
	if _, dup := plugins[plugin.ServerType][name]; dup {
		panic("plugin named " + name + " already registered for server type " + plugin.ServerType)
	}
	plugins[plugin.ServerType][name] = plugin
}

// EventName represents the name of an event used with event hooks.
type EventName string

// Define the event names for the startup and shutdown events
const (
	StartupEvent  EventName = "startup"
	ShutdownEvent EventName = "shutdown"
)

// EventHook is a type which holds information about a startup hook plugin.
type EventHook func(eventType EventName, eventInfo interface{}) error

// RegisterEventHook plugs in hook. All the hooks should register themselves
// and they must have a name.
func RegisterEventHook(name string, hook EventHook) {
	if name == "" {
		panic("event hook must have a name")
	}
	if _, dup := eventHooks[name]; dup {
		panic("hook named " + name + " already registered")
	}
	eventHooks[name] = hook
}

// EmitEvent executes the different hooks passing the EventType as an
// argument. This is a blocking function. Hook developers should
// use 'go' keyword if they don't want to block Corsair.
func EmitEvent(event EventName, info interface{}) {
	for name, hook := range eventHooks {
		err := hook(event, info)

		if err != nil {
			log.Printf("error on '%s' hook: %v", name, err)
		}
	}
}

// ParsingCallback is a function that is called after
// a directive's setup functions have been executed
// for all the server blocks.
type ParsingCallback func(Context) error

// RegisterParsingCallback registers callback to be called after
// executing the directive afterDir for server type serverType.
func RegisterParsingCallback(serverType, afterDir string, callback ParsingCallback) {
	if _, ok := parsingCallbacks[serverType]; !ok {
		parsingCallbacks[serverType] = make(map[string][]ParsingCallback)
	}
	parsingCallbacks[serverType][afterDir] = append(parsingCallbacks[serverType][afterDir], callback)
}

// SetupFunc is used to set up a plugin, or in other words,
// execute a directive. It will be called once per key for
// each server block it appears in.
type SetupFunc func(c *Controller) error

// DirectiveAction gets the action for directive dir of
// server type serverType.
func DirectiveAction(serverType, dir string) (SetupFunc, error) {
	if stypePlugins, ok := plugins[serverType]; ok {
		if plugin, ok := stypePlugins[dir]; ok {
			return plugin.Action, nil
		}
	}
	if genericPlugins, ok := plugins[""]; ok {
		if plugin, ok := genericPlugins[dir]; ok {
			return plugin.Action, nil
		}
	}
	return nil, fmt.Errorf("no action found for directive '%s' with server type '%s' (missing a plugin?)",
		dir, serverType)
}

// Loader is a type that can load a Corsairfile.
// It is passed the name of the server type.
// It returns an error only if something went
// wrong, not simply if there is no Corsairfile
// for this loader to load.
//
// A Loader should only load the Corsairfile if
// a certain condition or requirement is met,
// as returning a non-nil Input value along with
// another Loader will result in an error.
// In other words, loading the Corsairfile must
// be deliberate & deterministic, not haphazard.
//
// The exception is the default Corsairfile loader,
// which will be called only if no other Corsairfile
// loaders return a non-nil Input. The default
// loader may always return an Input value.
type Loader interface {
	Load(serverType string) (Input, error)
}

// LoaderFunc is a convenience type similar to http.HandlerFunc
// that allows you to use a plain function as a Load() method.
type LoaderFunc func(serverType string) (Input, error)

// Load loads a Corsairfile.
func (lf LoaderFunc) Load(serverType string) (Input, error) {
	return lf(serverType)
}

// RegisterCorsairfileLoader registers loader named name.
func RegisterCorsairfileLoader(name string, loader Loader) {
	corsairfileLoaders = append(corsairfileLoaders, corsairfileLoader{name: name, loader: loader})
}

// SetDefaultCorsairfileLoader registers loader by name
// as the default Corsairfile loader if no others produce
// a Corsairfile. If another Corsairfile loader has already
// been set as the default, this replaces it.
//
// Do not call RegisterCorsairfileLoader on the same
// loader; that would be redundant.
func SetDefaultCorsairfileLoader(name string, loader Loader) {
	defaultCorsairfileLoader = corsairfileLoader{name: name, loader: loader}
}

// loadCorsairfileInput iterates the registered Corsairfile loaders
// and, if needed, calls the default loader, to load a Corsairfile.
// It is an error if any of the loaders return an error or if
// more than one loader returns a Corsairfile.
func loadCorsairfileInput(serverType string) (Input, error) {
	var loadedBy string
	var corsairfileToUse Input
	for _, l := range corsairfileLoaders {
		cdyfile, err := l.loader.Load(serverType)
		if err != nil {
			return nil, fmt.Errorf("loading Corsairfile via %s: %v", l.name, err)
		}
		if cdyfile != nil {
			if corsairfileToUse != nil {
				return nil, fmt.Errorf("Corsairfile loaded multiple times; first by %s, then by %s", loadedBy, l.name)
			}
			loaderUsed = l
			corsairfileToUse = cdyfile
			loadedBy = l.name
		}
	}
	if corsairfileToUse == nil && defaultCorsairfileLoader.loader != nil {
		cdyfile, err := defaultCorsairfileLoader.loader.Load(serverType)
		if err != nil {
			return nil, err
		}
		if cdyfile != nil {
			loaderUsed = defaultCorsairfileLoader
			corsairfileToUse = cdyfile
		}
	}
	return corsairfileToUse, nil
}

// corsairfileLoader pairs the name of a loader to the loader.
type corsairfileLoader struct {
	name   string
	loader Loader
}

var (
	defaultCorsairfileLoader corsairfileLoader // the default loader if all else fail
	loaderUsed             corsairfileLoader // the loader that was used (relevant for reloads)
)
