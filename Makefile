default:
	cd cmd/abigen && go install && cd ../../
	cd cmd/crawler && go install

build-abigen:
	cd cmd/abigen && go build

build-resolver:
	cd cmd/resolver && go build listener.go resolver.go main.go

