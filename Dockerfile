#FROM golang
FROM golang:latest as builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
ENV GO111MODULE=on

WORKDIR /src
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download


RUN go build -o /app .
#RUN go install mytest .

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main .

EXPOSE 8000

ENTRYPOINT ["/app"]
CMD ["myapi"]