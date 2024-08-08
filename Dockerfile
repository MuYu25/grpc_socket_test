FROM golang:alpine as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

    #设置时区参数
ENV TZ=Asia/Shanghai
RUN sed -i 's!https://mirrors.ustc.edu.cn/!http://dl-cdn.alpinelinux.org/!g' /etc/apk/repositories

# RUN apk update --no-cache && apk add --no-cache tzdata
# RUN apk add --no-cache gcc
USER root
RUN apk cache clean
RUN apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

RUN apk --no-cache add gcc g++ make
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["make"]