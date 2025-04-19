## build. common build. 多阶段构建
FROM golang:1.24.2-alpine AS builder

WORKDIR /demo-server  
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

COPY api  ./api 
COPY dal ./dal
COPY model ./model
COPY service ./service
COPY cmd ./cmd
COPY router ./router
COPY handler ./handler
COPY docs ./docs
COPY conf  ./conf
COPY  go.mod go.sum ./ 

ENV GOPROXY=https://goproxy.io,direct
RUN go mod tidy 

# compile parameter.(by arg pass env, 可以从 docker-compose.yml文件传入)
ARG APP_ENV=prod
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /demo-server/biz-server ./cmd/demo-server


### production 构建
FROM alpine:latest AS prod 
WORKDIR /demo-server 
ARG APP_ENV=prod

# create log dir
RUN mkdir -p /demo-server/log /demo-server/conf

# copy builder阶段构建的结果 result.
COPY --from=builder /demo-server/biz-server .
# COPY conf/demo-server-${APP_ENV}.yaml ./conf/demo-server.yaml 在 prod 阶段，避免将构建时的配置文件复制到容器中，完全依赖运行时挂载

##  设置数据卷声明
VOLUME /demo-server/log 
VOLUME /demo-server/conf

ENTRYPOINT ["/demo-server/biz-server", "-f", "/demo-server/conf/demo-server-prod.yaml"]



### 开发环境的构建
FROM  golang:1.24.2-alpine AS dev 
WORKDIR /demo-server 
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN go mod tidy 

CMD ["go", "run", "./cmd/demo-server/main.go", "-f", "./conf/demo-server-dev.yaml"]

