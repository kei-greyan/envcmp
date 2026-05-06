// Package watcher provides a lightweight, poll-based file watcher for
// monitoring pairs of .env files.
//
// It is intentionally dependency-free and uses os.Stat to detect
// modification-time changes on a configurable interval. When either file
// changes, the caller-supplied onChange callback is invoked synchronously
// within the polling goroutine.
//
// Usage:
//
//	done := make(chan struct{})
//	err := watcher.Watch("left.env", "right.env", watcher.DefaultOptions(), func() {
//		// re-run comparison
//	}, done)
package watcher
