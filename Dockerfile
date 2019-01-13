FROM golang:1.11.4 as builder
LABEL authors="truegrom@gmail.com, aleksl0l@yandex.ru"
WORKDIR ${GOPATH}/src/books
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main

FROM scratch
LABEL authors="truegrom@gmail.com, aleksl0l@yandex.ru"
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]