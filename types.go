package main

type (
	PackageType int
)

const (
	type_WorldSave    PackageType = 1 // World save
	type_ResourcePack PackageType = 3 // Resource pack
	type_ModPack      PackageType = 5 // Mod pack
	type_Mod          PackageType = 6 // Mod
)
