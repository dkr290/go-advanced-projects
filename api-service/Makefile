
K3D_DEV_NAME    := "go-service-dev"
K3D_PRD_NAME    := "go-service-prd"
SALES_APP       := sales-api
BASE_IMAGE_NAME := localhost/api-service
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
DEPLOYMENT      := "apiservice-sales-api" 
SERVICE_NAME    := "SALES-API" 

run-local:
	go run ./app/services/sales-api/main.go  | go run ./app/tooling/logfmt/main.go -service=${SERVICE_NAME}

run-local-help:
	go run ./app/services/sales-api/main.go --help

tidy:
	go mod tidy
	go mod vendor


dev-create:
	k3d registry create; \
	k3d cluster create $(K3D_DEV_NAME) \
	  --api-port 6550 \
	  --servers 1 \
	  --agents 3 \
	  --k3s-arg --disable="traefik@server:0" \
	  --k3s-arg --disable="servicelb@server:0" \
	  --no-lb \
		--registry-use k3d-registry:37801 \
	  --wait

dev-start:
	k3d cluster start $(K3D_DEV_NAME)

dev-stop:
	k3d cluster stop $(K3D_DEV_NAME)

dev-delete:
	k3d cluster delete $(K3D_DEV_NAME); \
	k3d registry delete k3d-registry

dev-registry:
	k3d registry list

###---------DOCKER Build containers ---------------------------
all: sales

sales:
	docker build \
		-f zarf/docker/Dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# kubectl commands 
dev-restart:
	kubectl rollout restart deployment $(DEPLOYMENT) -n apiservice  

dev-logs:
	kubectl logs -n apiservice -l app.kubernetes.io/name=sales-api --all-containers=true -f --tail=100 | go run ./app/tooling/logfmt/main.go -service=${SERVICE_NAME}



