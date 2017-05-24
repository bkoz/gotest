# gotest
golang testing with OpenShift

```
oc new-app docker.io/jorgemoralespou/s2i-go~https://github.com/bkoz/gotest.git

oc create -f hpa.yaml

oc expose svc gotest

ab -n 1024 -c 4 http://<route>/
```


