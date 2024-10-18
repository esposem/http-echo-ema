FROM golang

RUN mkdir -p /app
RUN mkdir -p /app-attestation

# copy all required files: script, encrypted workload and enc/dec program
COPY http-echo.enc /app-attestation
COPY fenc /app-attestation
COPY http-echo /app
COPY run_server.sh /

EXPOSE 8080 8888
ENTRYPOINT ["/run_server.sh"]