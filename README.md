# kube-dash

Dashboard for sandbox environments

## Idea

The idea is to give the ability to start up services in a sandbox clusters
without giving them a lot of permissions (or any, for that mather).

![screenshot ](https://github.com/caarlos0/kube-dash/raw/master/screenshot.png)

## API

* `GET /api/deployments`: lists the all deployments in all namespaces except `kube-system`;
* `PUT /api/deployments/{namespace}/{deployment name}/up`: scale up a deployment (1 replica);
* `PUT /api/deployments/{namespace}/{deployment name}/down`: scale down a deployment (0 replicas).

## Deploy

Just run:

```sh
kubectl create -f https://raw.githubusercontent.com/caarlos0/kube-dash/master/deployment.yaml
```

It will create a deployment and a service in the `kube-system` namespace.

You will probably also need an ingress, here is an example:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: kube-dash
  namespace: kube-system
spec:
  rules:
  - host: dash.foo.local
    http:
      paths:
      - backend:
          serviceName: kube-dash
          servicePort: 80
```

