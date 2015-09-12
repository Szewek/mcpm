package database

type Database interface {
	func Read()
	func Update()
}