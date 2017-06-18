# Gravity Build System
(Components: Fusion, ...)
Collaborative, Compartmentalized, Build System using full virtualization to provide an environment that is reproducible, capable of colalborative editing and automates all the mundane aspects of development enabling greater collaboration of open source projects and enable new developers to join quickly without needing to deal with dependencies or other issues. just run the VM or participate on a system on someone elses infrastructure

## TODO:
### General
* Need to move all configuration into the gravity.build.yaml
* Switch to using corsiar as the base, then just pull out the build parts out of corsair
* Use the search plugin over and implement the current template structure
### WebUI
* Add Bolt, stick files cached into Bolt DB
* Use bleve, a elastic search type search for BoltDB
* Add globing to select groups of files, smart folder type system
* Read project *.YAML file, changes to file update UI, changes to webUI update file
* Watch folders and run action based on smart folder type system
* Collaborative editor
* Workviews, an open session in the project, either by HTML5 Xpra, or websockets+qemu, collaborative, but also multiple instances so different parts can be worked on at the same time
* Todo/Issues list that will be used to define version numbers for automatic versioning
* Release methods, scripts and pipleline of events or scripts to release software every version change
* webshell that interacts with a Qemu instance that runs the project

### Todo (different styles scrum, flat+tags,etc)
* Add hooks on completion of all with tag
* Add hooks on completuion of 1

### Virtual Filesystem
* Create a virtaul filesystem that can be mounted by a Qemu VM or used in other ways
* Integrity checking, versioning (using git)
* Easy rollback (right click), simplifying the git interface
* Merkel tree for torrent transfer

### Commands / Scheduling / Hooking
* Setup commands, for hooking onto events or other things
* chain commands into pipelines
* generate command line commands with subcommands like rake/make/etc
* schedule jobs

### QEMU VM
* Allow mock inputs into project
* ALlow mock outputs (mock APIs, etc) for output form project
* Tests, run after features are added (based on todo/issues maybe)
* Web shell and CLI shell to interact with the program
* Support multi server environments for complex testing and mocking and ensure system works from scratch

### Git
* Manage several projects, automated pushing based on finishing issue/item from todo

### TextUI
* shell to interact with gravity
* 

