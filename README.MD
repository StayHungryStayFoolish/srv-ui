# READ ME

1. run this command, generate go file in `proto`

         cd proto   
         protoc --go_out=plugins=grpc:. ****.proto

2. Generate gRPC stub 

        protoc -I/usr/local/include -I. \
          -I$GOPATH/src \
          -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
          --go_out=plugins=grpc:. \
          ****.proto
 
3. Generate reverse-proxy                  
        
        protoc -I/usr/local/include -I. \
          -I$GOPATH/src \
          -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
          --grpc-gateway_out=logtostderr=true:. \
          ***.proto
