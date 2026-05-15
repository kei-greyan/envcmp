// Package envpin provides snapshotting and drift detection for env maps.
//
// A "pin" records the expected keys, their presence, emptiness, and inferred
// value types (bool, int, float, string) at a point in time. The pin can be
// saved to a JSON file and later loaded to verify that a live env map has not
// drifted — i.e., that no keys have gone missing, become empty, or changed
// their inferred type.
//
// Typical usage:
//
//	pin := envpin.Create(env)
//	envpin.SaveFile(".env.pin", pin)
//
//	// later…
//	pin, _ := envpin.LoadFile(".env.pin")
//	violations := envpin.Check(pin, liveEnv)
package envpin
