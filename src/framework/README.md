## Corsair Application Framework
**Corsair is a system application framework, designed to provide very simple, lightweight structure for system applications.**

### Components
* **Configuration** - A basic framework for providing configurations for your system application, with default pathing based on modern POSIX defaults `$CORSAIR_USER/local/share/$CORSAIR_APPLICATION_NAME` and the `$CORSAIR_USER/config/$CORSAIR_APPLICATION_NAME`. Multiverse OS minimum design specifications presciribe avoiding any developer decisions that force developer preferences on users. In other words, support the ability to read and write all common formats, YAML, XML, JSON and customizable others. This allows the default to be selected by the user, since configurations should be stored in a secure way in memory and not waste available file descriptors or waste io, which is often the first bottlneck. 
* **Signal Handling** Intelligent signal handling, with easy to use defaults, and easy to configure. 
* **Init Script** Compile an init script and an install script that will make it very easy to install the application as an unpriviledged user, having this functionality by default will encourage novice users to learn about these systems and correct methods of running software, and prevent software from being run as root as often. 
* **CLI Flag Interface**
* **CLI Console Interface**
* **SSE and CSS only frameworks based WebUIs, rendered in a system window using a very simple HTML/CSS parser.**
* **Database handling with events to hook functionality**
* **Routing that will handle packets, HTTP requests, proxying and IPC**
* **Ephemeral key based authentication, directly using the design proposed for Multiverse OS**

A general MVC structure should be laid out so that it is very easy to see what to edit and make it trivial to get a hello world application that installs unprivledged easily, and comes with its own init script. Making it much easier to properly implment any system application from a simple folder watcher or complex p2p network client. 


