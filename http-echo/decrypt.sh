#! /bin/sh

cd /app

echo "#### Content of /app ####"
echo "+ ls /app"
ls /app

echo

if [ "$KATACLASS" = "cloud" ]; then
	echo "#### Fetching the key from trustee... ####"
	echo "+ curl http://127.0.0.1:8006/cdh/resource/default/workload-key/key.bin -o key.bin"
	curl http://127.0.0.1:8006/cdh/resource/default/workload-key/key.bin -o key.bin
else
	echo "#### Fetching the key from local cluster... ####"
	echo "+ curl https://people.redhat.com/eesposit/key.bin -o key.bin"
	curl https://people.redhat.com/eesposit/key.bin -o key.bin
fi

echo

echo "#### Content of /app ####"
echo "+ ls /app"
ls /app

echo

echo "#### Decrypting the workload ####"
echo "+ ./fenc -file http-echo.enc -key key.bin -operation decryption"
./fenc -file http-echo.enc -key key.bin -operation decryption

echo

echo "#### Running the workload ####"
chmod +x http-echo.enc.dec
echo "./http-echo.enc.dec"

./http-echo.enc.dec

