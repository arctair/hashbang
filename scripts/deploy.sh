VERSION=$1
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hashbang
  labels:
    app: hashbang
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hashbang
  template:
    metadata:
      labels:
        app: hashbang
    spec:
      containers:
      - name: hashbang
        image: arctair/hashbang:$VERSION
        ports:
        - containerPort: 5000
---
apiVersion: v1
kind: Service
metadata:
  name: hashbang
spec:
  type: NodePort
  ports:
  - port: 8080
    targetPort: 5000
  selector:
    app: hashbang
EOF
