module github.com/timickb/payments-example

go 1.22.3

replace github.com/timickb/narration-engine => ../../../

require (
	github.com/google/uuid v1.6.0
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/timickb/narration-engine v0.0.0-20240511230245-ccedc348f689
	google.golang.org/grpc v1.63.2
)

require (
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240429193739-8cf5692501f6 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
