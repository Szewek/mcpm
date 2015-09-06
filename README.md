# mcpm – Minecraft Package Manager
[![Build Status](https://travis-ci.org/Szewek/mcpm.svg?branch=master)](https://travis-ci.org/Szewek/mcpm)

`mcpm` is a package manager that lets you manage mods, saves and resource packs. It's written in Go.

This repository IS NOT related to [mcpm/mcpm](https://github.com/mcpm/mcpm).

## How does this work?
`mcpm` uses Curse CDN to gather information about everything Minecraft-related. So everything downloadable in Curse is also available in `mcpm`.

Package name comes from its curse.com URL. For example: /mc-mods/minecraft/**`tinkers-construct`**

## Installation (requires [Go](https://golang.org/) to compile)
```
go get github.com/Szewek/mcpm
```

Make sure you have set GOPATH/bin in PATH environment variable.

For first time you need to update database.
```
mcpm update
```

## Use examples
Command | Status | Description
--------|--------|------------
`mcpm get tinkers-construct` | Working | Downloads the newest version of [Tinkers' Construct](http://www.curse.com/mc-mods/minecraft/tinkers-construct) and puts in "mods" folder where this command was executed
`mcpm get tinkers-construct -for 1.7.10` | Not implemented | Downloads the latest version of that mod for Minecraft 1.7.10
`mcpm get tinkers-construct -d` | Not implemented | Only downloads that mod (does not put into subfolder)
`mcpm search Tinkers` | Working | Searches database for packages containing word "Tinkers" in package name, title and description
`mcpm update` | Working | Updates database
`mcpm forge` | Not implemented | Installs Minecraft Forge (recommended version)
`mcpm forge --latest` | Not implemented | Installs the latest version of Minecraft Forge
`mcpm authors` | Not implemented | Shows a list of all mcpm contributors
`mcpm authors tinkers-construct` | Not implemented | Shows mod authors
`mcpm make-server` | Not implemented | Downloads and installs Minecraft server instance (with Forge)
`mcpm list` | Not implemented | Lists all package names

## Contributing
You can submit bugs and requests. You are also allowed to modify this source code (fork it first, then create pull request).

## To do
- [x] Caching database
- [x] Getting package by unique name
- [ ] Getting package for appropriate version of Minecraft
- [ ] Unpacking modpacks
- [ ] Creating modpacks
- [ ] Creating server and client instances
- [ ] Getting mods' source code
