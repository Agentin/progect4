# Практика 4
## Выполнил: Студент ЭФМО-02-25 Фомичев Александр Сергеевич
### Структура:
```
services
    deploy
        monitoring
            prometheus.yml
            docker-compose.yml
    auth
        cmd
            auth
                main.go
        internal
            grpc
                server.go
            http
                handlers
                    login.go
                    verify.go
                routes.go
            service
                auth.go
    tasks
        cmd
            tasks
                main.go
        internal
            metrics
                metrics.go
            grpcclient
                client.go
            http
                middleware
                    metrics.go
                handlers
                    tasks.go
                    middleware
                        auth.go
                routes.go
            service
                tasks.go
shared
    shared
        logger
            logger.go 
    middleware
        requestid.go
        accesslog.go
        grpclog.go
    httpx
        client.go
pkg
    api
        auth
            v1
                auth.proto
                auth.pb.go
                auth_grpc.pb.go
docs
    pz17_api.md
README.md
go.mod
go.sum
```
### описание метрик
http_requests_total	Counter	method, route, status	Общее количество запросов
http_request_duration_seconds	Histogram	method, route	Длительность запросов в секундах
http_in_flight_requests	Gauge	(нет)	Текущее число одновременно выполняемых запросов
### Пример вывода /metrics 

**работа task**

![здесь должен быть рисунок, честно](image/4_2.png)

**работа task в Prometheus**

![здесь должен быть рисунок, честно](image/4_3.png)

**выводd /metrics**
![здесь должен быть рисунок, честно](image/4_1.png)

### Графики

**RPS**

![здесь должен быть рисунок, честно](image/4_4.png)

**error**

![здесь должен быть рисунок, честно](image/4_5.png)

**p95**

![здесь должен быть рисунок, честно](image/4_6.png)
