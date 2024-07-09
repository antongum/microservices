# Установка 
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- export GO_PATH=~/go
- export PATH=$PATH:/$GO_PATH/bin
- go get google.golang.org/grpc

# Генерация кода
- protoc -I=api --go_out=. movie.proto - - генерация protobuf 
- protoc -I=api --go_out=. --go-grpc_out=. movie.proto - генерация protobuf и grpc
- -I=api - положить в папку api в корне
- --go_out=. movie.proto - в какой файл положить protobuf (папка указана в файле .proto: option go_package = "/gen";)
- --go-grpc_out=. - в какой файл положить grpc (папка указана в файле .proto: option go_package = "/gen";)