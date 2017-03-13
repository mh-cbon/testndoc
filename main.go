// Package testndoc generates API documentation by listening to your tests.
package testndoc

// Recorder is a singleton to record api requests/responses.
var Recorder = &APIRecorder{}
