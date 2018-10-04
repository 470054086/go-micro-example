FROM docker.io/yam8511/go-micro
ADD ./ /go/src/greeter/api
WORKDIR /go/src/greeter/api
CMD micro  --registry_address 192.168.38.128:8500 api

