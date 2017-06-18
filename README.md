# Corsair
**A light weight server-focused application framework that is capable of routing based on packets, hostname, port, and protocol by default, based on a fork of Caddy**

Corsair is a lightweight, modular, application framework for Go, that by default provides a web server capable of automatically obtaining Let's Encrypt keys for every registered endpoint, and handling packets directly reading from the device and storing them in a ring buffer. Packets can then be routed to external or local endpoints, or handled by middleware plugins included at compile time. 

The developmentt will refocus on security, privacy, and making it easier to build and compile software ontop of corsair as plugins and deploying applications with a simple toolset for installing init scripts, and deploying software with debian packaging. 

The default web server plugin will be modified extensively, by default the web server plugin will spoof server tokens from other web servers including providing default 50x responses using the default templates of the spoofed server (Apache, NGINX). The other major change planned for the web server plugin is moving away from listening on specific ports and instead handling connections at a packet level, enabling software defined networking providing a interface for filtering, intercepting, routing packets, or responding. All functionality beyond the most basic requirements of the application framework will be implented as a plugins so the framework remains usuable for a wide variety of purposes. Providing a consistent framework to build Go software, handling the common tasks of CLI flags, YAML configuration, interfacing with in-memory and persistent databases, console UI, web UI (and tools for deploying webUI using custom rendering tools rather than using webkit or chromium engines), process managment, signal management, init script creation, and possibly debian packaging. 

*Corsair is reset back to v0.1.0 and is not ready for general use, although some developers may find it useful. It is being developed as a part of the Multiverse OS Go development library. If you need a production ready server, you will find greater success with [Caddy](https://github.com/caddyserver/caddy), it has developed overtime in a quality production ready web server. Which makes it a great starting point for a Go application framework.

### Planned Changes
 * Corsair will focus on a design that facilitates its use as a library and the fork will see major changes to the configuration subsystem.
 * Corsair will by default handle packets directly and communicate with a device or set of devices, enabling very low level filtering, enabling improved DOS protection and security.
 * Additionally Corsair will focus on expanding the ability to function as a high quality proxy, VPN, and portable tcp/udp server. 
 * Debian package built into the build process
 * Dropping all legacy windows support, using windows to serve files on the internet is dangerous, and leads to making compromises and failing to optimize everything to posix environments. 
 * Rich signal handling functionality with easy default but complex customization
 * An IPC server to come default with the web server
 * Include a make/rake type structure, then build in the ability to construct an init script, package debian, and include all plugins within the /src/plugin folder in the build. Structure to build in a reproducible way, in an ephemeral and compartmentalized environment. Possibly a good fit for mixing with the Multiverse OS gravity build system.
* Intergrate the `gravity` build system to simplify the development process, make it easy to participate and custom compile.


### Changes Actualized
* Major structural changes

    src/plugins/* - Is where you find plugins now, which hopefully will be much similer to work with. 
    src/corsair/* - All the corsair library files.
    .             - The only files in the root directory now are related to building a binary, the library will work seamlessly, enabling the binary to function as an example and and the src/corsair/ folder to be dropped in to any Go project. 


### Why?
After reviewing the code of Caddy, I discovered that it was a great example project for a full-featured Go application, even including an install script that setups up an unpriviledged user. It successfully used the Middleware abstraction which is useful for many designs and way easy to build on an event driven system. I believe caddy itself can be used as a application framework but I found it lacked some functionality, and included extra weight in features that I did not find particularly useful by default. I also felt the build system was obtuse and difficult, encouraging many users to just go to the site and be advertised at, rather than be bothered building it, regardless if intentional or not. 

Keeping that all in mind, I decided to combine Caddy, or Corsair, the system application framework, with the build system I begun working on. Pairing these together provided all the features I used in almost every application, and I made them even easier to remove by restructuring the code to make it easier to use as a library, and consolidated all plugins to the import directory `plugins/$PACKAGE_NAME`.

There are also many other improvements and features that I desired, and it just made sense to fork the project and use it as part of an application framework with secure defaults, privacy enhanced defaults. The goal is to provide enough structure to get started immediately on the parts that make your application unique, providing easy to use and implement defaults, with rich customization hiding underneath. This is done by separating all the plugins from core functionality and from each other, so any plugin can be removed, but by default providing easy to make build scripts for compiling, testing and publishing, in addition handle configuration, signal handling, pid/process control, init script creation, and basic UI (CLI flags, CLI console, and WebUI using a simple HTML/CSS parser).

This application framework is being designed to build a small collection of Go application that will be paired with Rust applications to make up the Multiverse OS ecosystem, but we are trying to make it general purpose enough that any Go developer could pick it up and get started quickly, and avoid being stuck rebuilding their configuration setup, database setup, signal handling and so on each time they start even a simple script. Using the plugin strategy described above, it allows developers to take only the parts they want, implemented simply or heavily customized, so their software can be light weight as possible, but also feature all the benefits of a larger project. It is the skeleton to hang things on and is being designed to cater to a wide variety of projects.


### Development
Corsair is being developed as a library to be packaged with the Multiverse OS Go development libraries. Anyone is welcome to contribute, pull requests will receive code reviews and will be pulled if the pull meets the minimum specifications of the Multiverse OS design specifications. Support will be provided for all Multiverse OS projects, the easiest way to get help is opening an issue on the related repository.  


