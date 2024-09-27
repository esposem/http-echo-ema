FROM golang
WORKDIR /app

# copy all required files: script, encrypted workload and enc/dec program
COPY http-echo.enc /app
COPY fenc /app
COPY decrypt.sh /app

EXPOSE 8080 8888
ENTRYPOINT ["/app/decrypt.sh"]