### Install ingress controller

```sh
root@VM-12-4-ubuntu:~$ kubectl create -f ingress-nginx-deployment.yaml
namespace/ingress-nginx created
serviceaccount/ingress-nginx created
configmap/ingress-nginx-controller created
clusterrole.rbac.authorization.k8s.io/ingress-nginx created
clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx created
role.rbac.authorization.k8s.io/ingress-nginx created
rolebinding.rbac.authorization.k8s.io/ingress-nginx created
service/ingress-nginx-controller-admission created
service/ingress-nginx-controller created
deployment.apps/ingress-nginx-controller created
ingressclass.networking.k8s.io/nginx created
validatingwebhookconfiguration.admissionregistration.k8s.io/ingress-nginx-admission created
serviceaccount/ingress-nginx-admission created
clusterrole.rbac.authorization.k8s.io/ingress-nginx-admission created
clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
role.rbac.authorization.k8s.io/ingress-nginx-admission created
rolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
job.batch/ingress-nginx-admission-create created
job.batch/ingress-nginx-admission-patch created
```

### Generate key-cert

```sh
# use openssl (version >= 1.1.1f) on Linux, e.g. Ubuntu 20.04
# don't run on macOS, which is using LibreSSL
# instead, you can `brew install openssl` on macOS
root@VM-12-4-ubuntu:~$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=learncn.com/O=learncn" -addext "subjectAltName = DNS:learncn.com"
Generating a RSA private key
.............................................................................................+++++
........................................+++++
writing new private key to 'tls.key'
-----
```

### Create secret

```sh
kubectl create secret tls learncn-tls --cert=./tls.crt --key=./tls.key
```

### Create a ingress

```sh
kubectl create -f ingress.yaml
```

### Test the result

```sh
curl -H "Host: learncn.com" https://10.97.150.2 -v -k
```

### Acesss Services
```sh
root@VM-12-4-ubuntu:~$ kubectl get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.97.150.2     <none>        80:32183/TCP,443:31157/TCP   43h
ingress-nginx-controller-admission   ClusterIP   10.108.34.153   <none>        443/TCP                      43h

curl -H "Host: learncn.com" https://localhost:31157 -v -k #所访问的服务需要跟ingress.yaml创建的Ingress资源在同一namespace下
```

### Understand why ingress is not enough?

- tls: cypher, dhkey, TLSVersion
- header based L7 rule
- rewriting?
  - header rewriting
  - URI rewriting

### How rewrite is supported in nginx ingress

```
Name	Description	Values
nginx.ingress.kubernetes.io/rewrite-target	Target URI where the traffic must be redirected	string
nginx.ingress.kubernetes.io/ssl-redirect	Indicates if the location section is accessible SSL only (defaults to True when Ingress contains a Certificate)	bool
nginx.ingress.kubernetes.io/force-ssl-redirect	Forces the redirection to HTTPS even if the Ingress is not TLS Enabled	bool
nginx.ingress.kubernetes.io/app-root	Defines the Application Root that the Controller must redirect if it's in '/' context	string
nginx.ingress.kubernetes.io/use-regex	Indicates if the paths defined on an Ingress use regular expressions	bool
```
