# gotest
golang testing with OpenShift

As cluster admin, add limits and quotas to the project:
```
oc create -f limits.json  
oc create -f resource-quotas.yaml
```

As a regular user, create the app:
```
oc new-app docker.io/jorgemoralespou/s2i-go~https://github.com/bkoz/gotest.git

oc create -f hpa.yaml

oc expose svc gotest

oc get hpa -w
```

In a second terminal, busy up the app with requests (don't forget the trailing slash):
```
ab -n 1024 -c 4 http://<route>/
```


