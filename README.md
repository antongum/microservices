## 01 - Что здесь
- структураы папок
- архитектура repository(gateway), controller, handler
- тонкости реализации repository(gateway), controller, handler
- тонкости реализации mock repo

## 02 - Что здесь
- Service Discovery Consul
- Service Discovery в памяти для тестов 
[ЗАПУСК_CONSUL_В_DOCKER]
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
