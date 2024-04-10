FROM alpine

#Use the repository of Chinese Mainland
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \  
    apk update && \  
    apk add --no-cache bash curl gcc git go musl-dev

# Set GOPROXY environment variable
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /svc/
COPY . .
RUN go mod download && \  
    bash make.sh 

VOLUME /svc/runtime/
EXPOSE 8888
ENTRYPOINT ["./showta"]