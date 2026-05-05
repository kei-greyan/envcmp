// Package config defines the Config struct that carries all runtime options
// for an envcmp invocation, along with validation logic.
//
// A Config is typically constructed by the cmd/envcmp main package from
// parsed CLI flags and then passed down to the comparator, filter, and
// formatter layers.
package config
