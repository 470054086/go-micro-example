FROM docker.io/yam8511/go-micro
ADD ./ /go/src/greeter/srv
WORKDIR /go/src/greeter/srv
CMD go run main.go --registry_address 192.168.38.128:8500
#CMD go run main.go

#CMD go run main.go --registry_address 192.168.38.128:8500