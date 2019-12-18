
all: muxacct
always:
	@:

muxedaccount.go: muxedaccount.x
	GOPRIVATE='*' PATH=$${PATH}:$$(go env GOPATH)/bin go generate

muxacct: muxedaccount.go always
	go build

clean:
	go clean
	rm -f *~ muxedaccount.go

.PHONY: all always clean
