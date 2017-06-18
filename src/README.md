# Corsair Framework structure
I believe a strategy of removing the remote dependencies in favor of good defaults that can be overridden or customized locally, breaking up the src folder into distinct categories:

    * `framework` - This would replace the current corsair folder, what is inside the corsair folder should be broken up across a consistent structure or at least symbolically linked to the $CORSAIR_ROOT directory, to make it easy to maintain by making the structure obvious and familiar to structures found in other frameworks.
    * `build-system` - This folder will provide the tools necessary to watch folers for automatic rebuilding and restarting, ability to create macros/commands and assign them to gravity <commands> <subcommands>, including defaults that make sense in the context of the Corsair application framework. 
    * `plugins` = This folder will be a place to put active plugins, ideally it will eventually be automatically imported through creation of an import file based on the contents of the folder before building, or dynamically importing at runtime. 

The framework folder may be better to break up into at least models views and controllers, because this is a known abstraction and works well for general application development as long as it is easily overriden. If done right, it could dramatically speed up production without limiting the developer by having deep customization and easy removal of any existing feature. This would be accomplished by implementing *everything* in the framework as a plugin, including the models, views and controllers logical structure. 

## Why is the build system a separate entity from the framework?

It is common to seperate these logical abstractions, for example, in Rails, the software rake, the Ruby equivilent to make, is the build tool packaged with Rails, and is packaged with defaults that provide basic functionality for simplifying starting/stopping/migrating and other tasks.

Similary Gravity provides this functionality and more. Beyond just macros/command assignment like make, it can also launch these commands or customized versions on hooks, and hooks can be tied to a number of things including file changes, or custom events. The entire Corsair system will be event driven, so these events will provide surface for these or other hooks available to developers to simplify their unique workflow, while providing simple and easy to use defaults.

