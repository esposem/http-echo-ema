#! /bin/sh

# sleep infinity

if [ "$ATTESTATION" = "1" ]; then
    cd /app-attestation

    echo "#### Content of /app-attestation ####"
    echo "+ ls /app-attestation"
    ls /app-attestation

    echo

    echo "#### Fetching the key from trustee... ####"
    echo "+ curl http://127.0.0.1:8006/cdh/resource/default/workload-key/key.bin -o key.bin"
    curl http://127.0.0.1:8006/cdh/resource/default/workload-key/key.bin -o key.bin

    echo

    echo "#### Decrypting the workload ####"
    echo "+ ./fenc -file http-echo.enc -key key.bin -operation decryption"
    ./fenc -file http-echo.enc -key key.bin -operation decryption

    echo

    chmod +x http-echo.enc.dec
    echo "./http-echo.enc.dec"
    ./http-echo.enc.dec
else
    cd /app

    echo "#### Content of /app ####"
    echo "+ ls /app"
    ls /app

    echo

    echo "./http-echo"
    ./http-echo
fi

