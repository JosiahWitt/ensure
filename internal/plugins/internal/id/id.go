// Package id provides constants used throughout the plugins. Most constant strings match their name, and are used to prevent typos from sneaking in.
package id

const (
	EnsuringE = "ensuring.E"

	Mocks      = "Mocks"
	SetupMocks = "SetupMocks"
	Subject    = "Subject"

	NEW = "NEW"

	Ensure              = "ensure"
	Ignore              = "-"
	ExampleIgnore       = "`ensure:\"-\"`"
	IgnoreUnused        = "ignoreunused"
	ExampleIgnoreUnused = "`ensure:\"ignoreunused\"`"
)
