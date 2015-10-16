// Package modes contains all subcommands (called "modes").
//
// MCPM starts mode by its name given in a command line.
//
// To create a mode. Create new Go source file in this package directory and follow that example:
//  func anyname(mo *ModeOptions) {
//      // Do something there
//  }
//  func init() {
//      registerMode("command-name", anyname);
//  }
package modes
