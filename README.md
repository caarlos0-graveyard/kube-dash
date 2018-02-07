# kube-dash

Dashboard for sandbox environments

## Deploy

Just run:

```sh
kubectl create -f https://raw.githubusercontent.com/caarlos0/kube-dash/master/deployment.yaml
```

It will create a deployment and a service in the `kube-system` namespace.

You will probably also need an ingress:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: kube-dash
spec:
  rules:
  - host: dash.foo.local
    http:
      paths:
      - backend:
          serviceName: kube-dash
          servicePort: 80
```

