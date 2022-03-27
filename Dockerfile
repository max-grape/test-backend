ARG IMAGE_GO
ARG IMAGE_ALPINE

FROM $IMAGE_GO
ARG CWD
ARG GOOS
ARG GOARCH
WORKDIR $CWD
COPY . .
ENV GOFLAGS=-mod=vendor
RUN GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -v -o app

FROM $IMAGE_ALPINE
ARG CWD
WORKDIR /app
RUN apk add -U curl tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone && \
    apk del -v tzdata && \
    rm -rf /var/cache/apk/*
COPY --from=0 $CWD/app ./app
HEALTHCHECK --interval=10s --timeout=2s CMD ["curl", "http://127.0.0.1:8090/health", "| grep -w \"healthy\""]
CMD ["./app"]
