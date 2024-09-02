FROM golang:1.23 as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

    #设置时区参数
ENV TZ=Asia/Shanghai
# RUN sed -i 's!https://mirrors.ustc.edu.cn/!http://dl-cdn.alpinelinux.org/!g' /etc/apk/repositories

# RUN apk update --no-cache && apk add --no-cache tzdata
# RUN apk add --no-cache gcc
USER root
# RUN apk cache clean

# RUN apk --no-cache add gcc g++
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /work/
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
RUN chmod +x ./main
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置时区（以 Asia/Shanghai 为例）
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
ENV LANG C.UTF-8
EXPOSE 8080 8089 9001
CMD [ "./main" ]