# This project is abandoned. Check [MCTool](https://github.com/Szewek/mctool), a successor of MCPM.

# MCPM - Minecraft Package Manager
[![Build Status](https://travis-ci.org/Szewek/mcpm.svg?branch=master)](https://travis-ci.org/Szewek/mcpm)
[![GoDoc](https://godoc.org/github.com/Szewek/mcpm?status.svg)](https://godoc.org/github.com/Szewek/mcpm)
[![Go Report Card](https://goreportcard.com/badge/github.com/Szewek/mcpm)](https://goreportcard.com/report/github.com/Szewek/mcpm)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FSzewek%2Fmcpm.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FSzewek%2Fmcpm?ref=badge_shield)

MCPM is a package manager AND A LIBRARY, that lets you manage mods, modpacks, world saves and resource packs. It's all written in [Go](https://golang.org/).
You can use this as a server command-line tool or create whole new tool using MCPM as a library.

The idea of MCPM is creating easy and fully automated management of Minecraft resources.
It is dedicated to both servers and clients.

This project is still missing its target features. You can contribute by giving new ideas or bugs at Issues page or developing the source code.

## How does this work?
MCPM gets information of all Minecraft-related stuff available from Curse CDN server. So everything downloadable on Curse is also available in MCPM.

Package names come from their specific curse.com URL.
For example: /mc-mods/minecraft/**`tinkers-construct`**

## Why Go? / Why not Java?
The answer is simple. This is not a Minecraft mod, launcher nor a development library.
Also Go is simple and easy to understand language.

## Installation
Check [Releases](https://github.com/Szewek/mcpm/releases) for downloads

To get the latest build, get [Go](https://golang.org/) and type this command:
```
go get github.com/Szewek/mcpm
```

If you have it compiled, make sure you have set GOPATH/bin in PATH environment variable.

## Modes
Modes are commands available in MCPM. They can be used as plugins for your own tool.
To run a mode, type this following command:
```
mcpm <modename>
```

## Modes examples
Command | Description
--------|------------
`mcpm get tinkers-construct` | Downloads the newest version of [Tinkers' Construct](http://www.curse.com/mc-mods/minecraft/tinkers-construct) and puts in "mods" folder where this command was executed
`mcpm search Tinkers` | Searches database for packages containing word "Tinkers" in package name, title and description
`mcpm update` | Updates database
`mcpm info xyz` | Displays information about package "xyz"

## TO DO
- [x] Caching database
- [x] Getting package by unique name
- [ ] Getting package for appropriate version of Minecraft
- [x] Unpacking modpacks and world saves
- [ ] Creating server and client instances


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FSzewek%2Fmcpm.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FSzewek%2Fmcpm?ref=badge_large)
