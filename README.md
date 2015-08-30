# mcpm – Minecraft Package Manager
`mcpm` is a package manager for mods, saves and resource packs. It's written in Go.

This repository IS NOT related to [mcpm/mcpm](https://github.com/mcpm/mcpm).

## How this will work?
`mcpm` will use Curse CDN to gather information about everything Minecraft-related. So everything downloadable in Curse is also available in `mcpm`.

Package name comes from its curse.com URL. For example: /mc-mods/minecraft/**`tinkers-construct`**

## Examples
Command | Status | Description
--------|--------|------------
`mcpm get tinkers-construct` | Not implemented | Downloads the newest version of [Tinkers' Construct](http://www.curse.com/mc-mods/minecraft/tinkers-construct) and puts in "mods" folder where this command was executed
`mcpm get tinkers-construct --for 1.7.10` | Not implemented | Downloads the latest version of that mod for Minecraft 1.7.10
`mcpm get tinkers-construct -d` | Not implemented | Only downloads that mod
`mcpm search Tinkers` | Not implemented | Searches database for packages containing word "Tinkers" in package name, title and description
`mcpm update` | Not implemented | Updates database
`mcpm update -f` | Not implemented | Forces to update whole database
`mcpm forge` | Not implemented | Installs Minecraft Forge (recommended version)
`mcpm forge --latest` | Not implemented | Installs the latest version of Minecraft Forge

## Contributing
You can submit bugs and requests. You are also allowed to modify this source code (fork it first, then create pull request).

## To do
- [ ] Caching database
- [ ] Getting package by unique name
- [ ] Getting package for appropriate version of Minecraft
- [ ] Creating modpacks
- [ ] Creating server and client instances
- [ ] Getting mods' source code
