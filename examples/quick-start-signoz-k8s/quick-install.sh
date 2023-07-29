# Installing Tracetest and Signoz on a Kubernetes cluster

k3d cluster create tracetest-signoz

cat << EOF > opentelemetry-collector-resources.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: collector-config
data:
  collector.yaml: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:

    processors:
      batch:
        timeout: 100ms

      # Data sources: traces
      probabilistic_sampler:
        hash_seed: 22
        sampling_percentage: 100

    exporters:
      # Output logger, used to check OTel Collector sanity
      logging:
        loglevel: debug

      # OTLP for Tracetest
      otlp/tracetest:
        endpoint: tracetest.tracetest.svc.cluster.local:4317
        tls:
          insecure: true
      # OTLP for Signoz
      otlp/signoz:
        endpoint: signoz-otel-collector.observability.svc.cluster.local:4317
        tls:
          insecure: true

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [probabilistic_sampler, batch]
          exporters: [otlp/signoz, otlp/tracetest, logging]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otel-collector
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: otel-collector
  template:
    metadata:
      labels:
        app.kubernetes.io/name: otel-collector
    spec:
      containers:
        - name: otelcol
          args:
            - --config=/conf/collector.yaml
          image: otel/opentelemetry-collector-contrib:0.67.0
          volumeMounts:
            - mountPath: /conf
              name: collector-config
          resources:
            requests:
              cpu: 250m
              memory: 512Mi
      volumes:
        - configMap:
            items:
              - key: collector.yaml
                path: collector.yaml
            name: collector-config
          name: collector-config
---
apiVersion: v1
kind: Service
metadata:
  name: otel-collector
spec:
  ports:
    - name: grpc-otlp
      port: 4317
      protocol: TCP
      targetPort: 4317
  selector:
    app.kubernetes.io/name: otel-collector
  type: ClusterIP
EOF

cat << EOF > tracetest-resources.yaml
---
# Source: tracetest/charts/postgresql/templates/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: tracetest-postgresql
  namespace: "tracetest"
  labels:
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-12.1.6
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/managed-by: Helm
type: Opaque
data:
  postgres-password: "eVJPeDZUcWZrcg=="
  password: "bm90LXNlY3VyZS1kYXRhYmFzZS1wYXNzd29yZA=="
  # We don't auto-generate LDAP password when it's not provided as we do for other passwords
---
# Source: tracetest/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tracetest
  labels:
    helm.sh/chart: tracetest-0.2.69
    app.kubernetes.io/name: tracetest
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/version: "v0.13.0"
    app.kubernetes.io/managed-by: Helm
data:
  config.yaml: |-
    postgres:
      host: tracetest-postgresql
      user: tracetest
      password: not-secure-database-password
      port: 5432
      params: sslmode=disable

    telemetry:
      exporters:
        collector:
          serviceName: tracetest
          sampling: 100 # 100%
          exporter:
            type: collector
            collector:
              endpoint: otel-collector.tracetest.svc.cluster.local:4317

    server:
      httpPort: 11633
      otlpGrpcPort: 4317
      otlpHttpPort: 4318
      telemetry:
        exporter: collector

  provisioning.yaml: |-
    type: DataStore
    spec:
      id: current
      name: Signoz
      type: signoz
---
# Source: tracetest/charts/postgresql/templates/primary/svc-headless.yaml
apiVersion: v1
kind: Service
metadata:
  name: tracetest-postgresql-hl
  namespace: "tracetest"
  labels:
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-12.1.6
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: primary
    # Use this annotation in addition to the actual publishNotReadyAddresses
    # field below because the annotation will stop being respected soon but the
    # field is broken in some versions of Kubernetes:
    # https://github.com/kubernetes/kubernetes/issues/58662
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  type: ClusterIP
  clusterIP: None
  # We want all pods in the StatefulSet to have their addresses published for
  # the sake of the other Postgresql pods even before they're ready, since they
  # have to be able to talk to each other in order to become ready.
  publishNotReadyAddresses: true
  ports:
    - name: tcp-postgresql
      port: 5432
      targetPort: tcp-postgresql
  selector:
    app.kubernetes.io/name: postgresql
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/component: primary
---
# Source: tracetest/charts/postgresql/templates/primary/svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: tracetest-postgresql
  namespace: "tracetest"
  labels:
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-12.1.6
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: primary
  annotations:
spec:
  type: ClusterIP
  sessionAffinity: None
  ports:
    - name: tcp-postgresql
      port: 5432
      targetPort: tcp-postgresql
      nodePort: null
  selector:
    app.kubernetes.io/name: postgresql
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/component: primary
---
# Source: tracetest/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: tracetest
  labels:
    helm.sh/chart: tracetest-0.2.69
    app.kubernetes.io/name: tracetest
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/version: "v0.13.0"
    app.kubernetes.io/managed-by: Helm
spec:
  type: ClusterIP
  ports:
    - port: 11633
      targetPort: http
      protocol: TCP
      name: http
    - port: 4317
      targetPort: otlp-grpc
      protocol: TCP
      name: otlp-grpc
    - port: 4318
      targetPort: otlp-http
      protocol: TCP
      name: otlp-http
  selector:
    app.kubernetes.io/name: tracetest
    app.kubernetes.io/instance: tracetest
---
# Source: tracetest/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tracetest
  labels:
    helm.sh/chart: tracetest-0.2.69
    app.kubernetes.io/name: tracetest
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/version: "v0.13.0"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: tracetest
      app.kubernetes.io/instance: tracetest
  template:
    metadata:
      labels:
        app.kubernetes.io/name: tracetest
        app.kubernetes.io/instance: tracetest
    spec:
      serviceAccountName: default
      securityContext:
        {}
      containers:
        - name: tracetest
          securityContext:
            {}
          image: "kubeshop/tracetest:v0.13.0"
          env:
          - name: TRACETEST_DEV
            value: "true"
          args:
          - --config
          - '/app/config/config.yaml'
          - --provisioning-file
          - '/app/config/provisioning.yaml'
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: 11633
              protocol: TCP
            - name: otlp-grpc
              containerPort: 4317
              protocol: TCP
            - name: otlp-http
              containerPort: 4318
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            requests:
              cpu: 250m
              memory: 512Mi
          volumeMounts:
          - name: config
            mountPath: /app/config
      volumes:
      - name: config
        configMap:
          name: tracetest
---
# Source: tracetest/charts/postgresql/templates/primary/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tracetest-postgresql
  namespace: "tracetest"
  labels:
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-12.1.6
    app.kubernetes.io/instance: tracetest
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: primary
  annotations:
spec:
  replicas: 1
  serviceName: tracetest-postgresql-hl
  updateStrategy:
    rollingUpdate: {}
    type: RollingUpdate
  selector:
    matchLabels:
      app.kubernetes.io/name: postgresql
      app.kubernetes.io/instance: tracetest
      app.kubernetes.io/component: primary
  template:
    metadata:
      name: tracetest-postgresql
      labels:
        app.kubernetes.io/name: postgresql
        helm.sh/chart: postgresql-12.1.6
        app.kubernetes.io/instance: tracetest
        app.kubernetes.io/managed-by: Helm
        app.kubernetes.io/component: primary
      annotations:
    spec:
      serviceAccountName: default

      affinity:
        podAffinity:

        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app.kubernetes.io/name: postgresql
                    app.kubernetes.io/instance: tracetest
                    app.kubernetes.io/component: primary
                topologyKey: kubernetes.io/hostname
              weight: 1
        nodeAffinity:

      securityContext:
        fsGroup: 1001
      hostNetwork: false
      hostIPC: false
      initContainers:
      containers:
        - name: postgresql
          image: docker.io/bitnami/postgresql:14.7.0-debian-11-r29
          imagePullPolicy: "IfNotPresent"
          securityContext:
            runAsUser: 1001
          env:
            - name: BITNAMI_DEBUG
              value: "false"
            - name: POSTGRESQL_PORT_NUMBER
              value: "5432"
            - name: POSTGRESQL_VOLUME_DIR
              value: "/bitnami/postgresql"
            - name: PGDATA
              value: "/bitnami/postgresql/data"
            # Authentication
            - name: POSTGRES_USER
              value: "tracetest"
            - name: POSTGRES_POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: tracetest-postgresql
                  key: postgres-password
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: tracetest-postgresql
                  key: password
            - name: POSTGRES_DB
              value: "tracetest"
            # Replication
            # Initdb
            # Standby
            # LDAP
            - name: POSTGRESQL_ENABLE_LDAP
              value: "no"
            # TLS
            - name: POSTGRESQL_ENABLE_TLS
              value: "no"
            # Audit
            - name: POSTGRESQL_LOG_HOSTNAME
              value: "false"
            - name: POSTGRESQL_LOG_CONNECTIONS
              value: "false"
            - name: POSTGRESQL_LOG_DISCONNECTIONS
              value: "false"
            - name: POSTGRESQL_PGAUDIT_LOG_CATALOG
              value: "off"
            # Others
            - name: POSTGRESQL_CLIENT_MIN_MESSAGES
              value: "error"
            - name: POSTGRESQL_SHARED_PRELOAD_LIBRARIES
              value: "pgaudit"
          ports:
            - name: tcp-postgresql
              containerPort: 5432
          livenessProbe:
            failureThreshold: 6
            initialDelaySeconds: 30
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
            exec:
              command:
                - /bin/sh
                - -c
                - exec pg_isready -U "tracetest" -d "dbname=tracetest" -h 127.0.0.1 -p 5432
          readinessProbe:
            failureThreshold: 6
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
            exec:
              command:
                - /bin/sh
                - -c
                - -e

                - |
                  exec pg_isready -U "tracetest" -d "dbname=tracetest" -h 127.0.0.1 -p 5432
                  [ -f /opt/bitnami/postgresql/tmp/.initialized ] || [ -f /bitnami/postgresql/.initialized ]
          resources:
            limits: {}
            requests:
              cpu: 250m
              memory: 256Mi
          volumeMounts:
            - name: dshm
              mountPath: /dev/shm
            - name: data
              mountPath: /bitnami/postgresql
      volumes:
        - name: dshm
          emptyDir:
            medium: Memory
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes:
          - "ReadWriteOnce"
        resources:
          requests:
            storage: "8Gi"
EOF

# add signoz repo and install from helm chart
echo
echo "Installing Signoz..."
echo

helm repo add signoz https://charts.signoz.io
helm repo update
helm install signoz signoz/signoz --namespace observability --create-namespace

echo
echo "Signoz installed"
echo

# install tracetest and opentelemetry-collector
echo
echo "Installing Tracetest and OTel Collector..."
echo

kubectl create namespace tracetest
kubectl -n tracetest apply -f opentelemetry-collector-resources.yaml
kubectl -n tracetest apply -f tracetest-resources.yaml

echo
echo "Tracetest and OTel Collector installed."
echo

# wait for 60 seconds
sleep 60

# port-forward
kubectl port-forward --namespace tracetest svc/tracetest 11633

# kubectl port-forward --namespace observability svc/signoz-frontend 3301:3301 & \
# echo "Press CTRL-C to stop port forwarding and exit the script"
# wait
