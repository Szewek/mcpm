package main

type (
	_PackageType int
)

const (
	type_WorldSave    _PackageType = 1 // World save
	type_ResourcePack _PackageType = 3 // Resource pack
	type_ModPack      _PackageType = 5 // Mod pack
	type_Mod          _PackageType = 6 // Mod
)
