# 初始微服务
- 本文使用go-micro框架进行微服务开发,本文的代码都是从 https://github.com/micro/examples/tree/master/greeter 稍微进行修改的
- 想了解微服务的定义 请查看本篇文章 https://wuyin.io/2018/05/10/microservices-part-1-introduction-and-consignment-service/ 本次的配置也很多是来源这篇文章


- 开始安装微服务
- 安装go-micro
  go get -u github.com/micro/go-micro
>   此次安装因为含有google的grpc和crypto包,肯定会有被强的可能,可以手动安装这两个包,然后放入gopath之中  https://github.com/golang/crypto,https://github.com/grpc/grpc-go
- 安装protobuf
  - go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
  - go get -u github.com/micro/protoc-gen-micro
>   此两个插件可以proto转化为各自的代码,go-micro使用的google的grpc,grpc使用的protobuf序列化协议
- 安装micro命令行工具
  - go get -u github.com/micro/micro
  - 安装好之后,生成二进制文件 go build GOPATH/src/github.com/micro/micro 生成之后,加入环境变量,方便使用
## 服务发现
微服务最重要的功能,就是服务注册和发现,go-micro依赖于consul(可以通过插件自行更换etcd,zoopker)实现服务发现,服务发现的介绍 https://www.jianshu.com/p/5eac16e9804a
- 安装consul
  - https://www.consul.io 下载consul
  - 启动consul  consul agent --dev

### 代码书写
![image](https://images2015.cnblogs.com/blog/15172/201612/15172-20161225120450964-1678545462.png)


1. 微服务的架构图大概如此,先启动service,service启动了进行consul进行注册,启动client,client通过consul发现服务,最外层api-getway用于接受客户请求,调取client,clinet进行服务发现,然后负载均衡的调取service端,consul的服务注册和发现是服务断开,自动从连接中去掉这个service,可以不通过配置,完全负载均衡。
2. 本次依次实现两个服务端,一个say-service,一个order-service,一个客服端api

### 代码大致
> 可以使用micor new example进行末班代码的生成


1.say-service的实现
```go

//先定义proto
//使用什么协议
syntax = "proto3";
###定义的包名
package go.micro.srv.greeter;
###定义服务
service Say {
    ##定义rpc方法
	rpc Hello(Request) returns (Response) {}
	rpc World(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}

###然后使用刚开始下载的工具 进行编译

protoc --proto_path=. --micro_out=. --go_out=. ./proto/saytwo/saytwo.proto

定义service启动 main.go


func main() {
    ### 定义服务
	service := micro.NewService(
	    //这个名字相当重要,用于和consul进行通信,唯一性
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	// optionally setup command line usage
	service.Init()
	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}


### 实现服务

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello,我是第一个 " + req.Name
	return nil
}

func (s *Say) World(ctx context.Context,req *hello.Request,res *hello.Response) error{
	log.Print("Received Say.World request")
	res.Msg = "World" + req.Name
	return nil
}

### 这里大概为代码的简略实现,具体可以参考github,type定义可以想象为php中的class func为class下的method
```

服务定义好了,现在先启动consul,然后在启动服务,访问consul的web会发现服务已经注册进去了 http://localhost:8500

2. order-service的实现大致相同
3. api-clinet的实现
```go
// proto定义的实现大致相同
//定义服务的启动

func main() {
	//定义服务名字
	service := micro.NewService(
		micro.Name("go.micro.api.greeter"),
	)

	// parse command line flags
	service.Init()
	//定义远程调用
	service.Server().Handle(
		service.Server().NewHandler(
		    //说明了服务的名字
			&say.Say{Client: hello.NewSayService("go.micro.srv.greeter", service.Client())},
			//&Order{Client:  sayto.NewOrderService("go.micro.srv.greeter", service.Client())},
		),
	)
	//定义使用路由 这里应该有简单方法 暂时未找到 先这样吧
	service.Server().Handle(
		service.Server().NewHandler(
			//&Say{Client: hello.NewSayService("go.micro.srv.greeter", service.Client())},
			&order.Order{Client: sayto.NewOrderService("go.micro.srvtwo.greeter", service.Client())},
		),
	)
	//启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```
启动api-clinet 我们会发现clinet也注入进了consul,达到目的

4. api-getway
使用我们刚刚下载的命令,自动会实现api-getway  **micro api ** 就实现了一个api-getway,并且也已经注册到consul

5. 测试微服务
   http://localhost/greeter/order/yes?name=John 返回了json数据,完美,其实他的api-getway转发有点类似php实现的根据module/controller/method进行路由转发的原理一直


#### 微服务docker化
> docker微服务最重要的就是docker之间的通信,这里我们使用了--registry_address 192.168.38.128:8500,即是使用的ip地址,还可以使用link,这里暂时不讨论

> docker启动的时候,请注意防火墙的原因,防火墙关闭了可能会导致端口无法分享,docker打开了可能会导致docker之间无法通信

1. 既然实现了微服务,不进行docker化又怎么是微服务呢？依据"一个docker即一个服务"的原则,进行docker话,我们的此微服务大概需要下面几个docker
```ini
1. consul docker
2. api-getway docker
3. api-clinet docker
4. say-service docker
5. order-service docker
```
2. 实现consul docker
```ini
docker pull consul
###启动consul,这个我们需要对外共享8500 8400 8300
docker run -d  -p 8300:8300 -p 8400:8400 -p 8500:8500 -p 8600:53/udp   docker.io/consul

```
3. 实现api-getway

```ini
###这个我们使用Dockerfile来实现构造docker
//Dockerfile
###基础镜像
FROM docker.io/yam8511/go-micro
###docker启动执行的命令
CMD micro  --registry_address host:8500 api

###构造docker 此处. 表示docker在当前目录

docker build -t micro-service/api-getway .

###启动docker api-getway需要对外分享端口

docker run -d -p 8080:8080 micro-service/api

```

4. 实现api-client
```ini


###这个我们使用Dockerfile来实现构造docker

### 来源哪里
FROM docker.io/yam8511/go-micro
### 将当前目录都添加至dest目录
ADD ./ /go/src/greeter/api
//docker的工作目录
WORKDIR /go/src/greeter/api
##执行命令
CMD go run api.go --registry_address 192.168.38.128:8500

### 构造docker

docker build -t micro-service/api-clinet .

###启动docker

docker run -d micro-service/api-client

```

5. 实现order-service
``` ini

FROM docker.io/yam8511/go-micro
ADD ./ /go/src/greeter/srvtwo
WORKDIR /go/src/greeter/srvtwo
CMD go run main.go --registry_address 192.168.38.128:8500

### 构造docker

docker build -t micro-service/order-service .

###启动docker

docker run -d micro-service/order-service

```
6. 实现say-service 同order-service一致

### docker-compose 进行docker依赖

> 大家可能也发现了上面的问题,对于每一个服务我都要进行一次dockerFile 其实相当麻烦,这里我们就使用docker-compose进行服务依赖



### 服务依赖

```
version: '3.0'
services:
  consul:
    image: consul
    hostname: "registry"
    ports:
    - "8300:8300"
    - "8400:8400"
    - "8500:8500"
    - "8600:53/udp"

  api:
#    command: --registry_address=registry:8500 --register_interval=5 --register_ttl=10 api
    image: micro-service/api
    links:
    - consul
    ports:
    - "8080:8080"
  api-service:
     image: micro-service/api-service
     links:
     - consul
     - api
  order-service:
     image: micro-service/order-service
     links:
     - consul
     - api
  say-service:
     image: micro-service/say-service
     links:
     - consul
     - api
```
docker-compose up -d

> 服务依赖启动之后,会存在有的请求访问不到,暂时未找到原因