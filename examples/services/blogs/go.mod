module github.com/timickb/blogs-example

go 1.22.3

//replace github.com/timickb/narration-engine => ../../../

require (
	github.com/google/uuid v1.6.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.59.0
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/timickb/narration-engine v0.0.0-20240516185045-1ef813193d0d // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
