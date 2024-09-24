FROM golang
WORKDIR /app
COPY . /app
RUN go build -tags netgo
EXPOSE 8080 8888
USER 1001
ENTRYPOINT ["/app/http-echo-ema"]