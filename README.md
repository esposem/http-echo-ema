## Usage

### Considerations:

1. route must be added to /etc/hosts: `<your ip>      my-web-app-route-default.apps.coco2410.kata.com`
2. key.bin must be uploaded in the Trustee. The key used in the current http server is in https://people.redhat.com/~eesposit/key.bin
3. Check the ip in `cloud.yaml` line 25
4. Check the env variables in both `cloud.yaml` and `prem.yaml`: `ATTESTATION` (only to be set in `cloud.yaml`) to try attestation and `DD_MB_SIZE` to decide how big each request should occupy in memory. Do not touch the other one.
5. Check pod work `APP_URL=$(oc get routes/my-web-app-route -o jsonpath='{.spec.host}')`:
 `curl $APP_URL` and notice how the ip address changes
	```
	# curl $APP_URL
	Service handled by pod on-prem-deployment-68fcf4d4f-qw4rj
	# curl $APP_URL
	Service handled by pod on-prem-deployment-68fcf4d4f-g5nff
	```
1. before running count_calls.sh you must have logged in with `oc`

### Run:

1. run your KEDA policy
2. modify the env vars and then run cloud.yaml
3. modify the env vars and then run prem.yaml
4. run service.yaml
5. run route.yaml
6. run count_calls.sh and observe the result

