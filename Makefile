default:
	cd cmd/ && go build && mv cmd ../execfile
run:
	./execfile --config $(CONFIG)

build-abigen:
	cd cmd/abigen && go build

build-resolver:
	cd cmd/resolver && go build listener.go resolver.go main.go

 
