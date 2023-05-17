# Setup a Kubernetes cluster with OTel enabled using K3D

1. Install k3d:

```sh
brew install k3d
#or
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```

2. Go to `k3d-setup` folder and start Jaeger and OTel Collector with `docker compose`:

```sh
cd ./k3d-setup
docker compose up -d
```

3. Run the following command, replacing the placeholder `PATH_TO_THIS_FOLDER` with the current folder absolute path (that can get with `echo $PWD`):
```sh
k3d cluster create tracingcluster \
  --image=rancher/k3s:v1.27.1-k3s1 \
  --volume '[PATH_TO_THIS_FOLDER]/containerd/config.toml.tmpl:/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl@server:*' \
  --volume '[PATH_TO_THIS_FOLDER]/kube-tracing:/etc/kube-tracing@server:*' \
  --k3s-arg '--kube-apiserver-arg=tracing-config-file=/etc/kube-tracing/apiserver-tracing.yaml@server:*' \
  --k3s-arg '--kube-apiserver-arg=feature-gates=APIServerTracing=true@server:*' \
  --k3s-arg '--kubelet-arg=config=/etc/kube-tracing/kubelet-tracing.yaml@server:*'
```

4. After setting up the cluster and Jaeger, you should be able to see Kubernetes traces with Jaeger on `http://localhost:16686`.

5. A single test that you can do is to run: 

```sh
kubectl -v 9 run -it --rm --restart=Never --image=alpine alpine -- echo hi
```
