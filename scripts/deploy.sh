VERSION=$1
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-starter
  labels:
    app: go-starter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-starter
  template:
    metadata:
      labels:
        app: go-starter
    spec:
      containers:
      - name: go-starter
        image: arctair/go-starter:$VERSION
        ports:
        - containerPort: 5000
---
apiVersion: v1
kind: Service
metadata:
  name: go-starter
spec:
  type: NodePort
  ports:
  - port: 8080
    targetPort: 5000
  selector:
    app: go-starter
EOF
