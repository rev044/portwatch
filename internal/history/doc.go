// Package history implements a bounded, thread-safe ring buffer for storing
// recent port change events observed by portwatch.
//
// Usage:
//
//	h := history.New(200)   // keep the last 200 events
//	h.Record(change)        // called from the alert/notifier pipeline
//	for _, e := range h.Entries() {
//		fmt.Println(e.Timestamp, e.Change)
//	}
//
// The ring buffer overwrites the oldest entry once the capacity is reached,
// keeping memory usage constant regardless of runtime duration.
package history
