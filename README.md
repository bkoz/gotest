# Autoscaling go applications with OpenShift

As cluster admin, add limits and quotas to the project:
```
oc create -f limits.json  
oc create -f resource-quotas.yaml
```

As a regular user, create the application, route and autoscaler objects.
```
oc new-app docker.io/jorgemoralespou/s2i-go~https://github.com/bkoz/gotest.git

oc create -f hpa.yaml

oc expose svc gotest --path=/mandelbrot

oc get hpa -w
```

In a second terminal, busy up the app with requests anmd wait for autoscaling to happen (don't forget the trailing slash):
```
ab -n 20000 -c 16 http://<route>/
```


