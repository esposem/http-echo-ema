## Usage

1. Add the service to /etc/hosts: `<your ip>      ema-http-echo-service-default.apps.coco2410.kata.com`
2. change `kyverno_policy.yaml` to have your system ip there
3. run `deployment.yaml`
4. run `service.yaml`
5. `oc expose service ema-http-echo-service -l app=ema-http-echo`
6. run `kyverno_policy.yaml`
7. scale deployment to two pods
8. check the pods have different runtimeclass
9. `APP_URL=$(oc get routes/ema-http-echo-service -o jsonpath='{.spec.host}')`
10. `curl $APP_URL` and notice how the ip address changes
	```
	# curl $APP_URL
	Hello from pod ema-http-echo-68fcf4d4f-qw4rj
	# curl $APP_URL
	Hello from pod ema-http-echo-68fcf4d4f-g5nff
	```