# Corsair Plugins
*Corsair is a networked system application framework, providing a simple framework with minimal complexity defaults, but with deep customization*

Consider the possibility of using runtime plugins, by reading all packages within a certain directory, compiling them or copying them to a segment of the compiled binary, or dynamically loading them:

**Dynamic Loading** [tools/go/loader](https://godoc.org/golang.org/x/tools/go/loader)
*(I did not have time to do a deep review, but I believe this is the correct functionality. I know a recent version of Go added plugins, but I have not reviewd those enough to determine if they are a correct fit. Plugins may need to provide additional functionality of established struct objects.*)

## Setup
One just needs to add your package to this folder, then include an import in the main file in the src of your corsair application, with the _ symbol as a prefix to the path to the plugins. If it is placed in the exact structure `$CORSAIR_ROOT/src/plugins`, then one can import their plugin using just `"plugins/$PLUGIN_NAME"`. This makes it much easier to customize and maintain plugins. 

Conversion of Caddy plugins will work for the time being but there is little liklihood that the APIs will stay consistent on a long enough time line. We will try to avoid unncessary change and create abstraction overlays where possible but the two projects have different fundamental goals. 

## Boilerplate 
A boilerplate providing a working example, and eventual skeleton for automated code generation, is provided by default and time will be invested in heavily documenting it, incorporating different examples, providing additional documentation and keeping it consistent with any framework changes.


## Dynamic Loading
Eventually it would be nice to just generate a source file based on all the packages in the plugin folder and include it automatically to simplify plugin install. There are a variety of ways to accomplish this and time to review each one will be needed before a decision can be made on the best approach for Corsair.


