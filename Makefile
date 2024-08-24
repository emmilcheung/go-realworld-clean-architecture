# ==============================================================================
# Docker compose commands

.PHONY: docker-local
docker-local:
	echo "Starting docker environment"
	docker-compose up --build

.PHONY: docker-down
docker-down:
	echo "Cleaning docker environment"
	docker-compose down

# ==============================================================================
# Main

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: build
build:
	go build ./cmd/api/main.go

# ==============================================================================
# Modules support

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache


# ==============================================================================
# k8s support

.PHONY: k8s-starting
k8s-apply:
# prepare config / file system setup
	kubectl apply -f ./k8s/postgresql-pv.yaml
	kubectl apply -f ./k8s/postgresql-pvc.yaml
	kubectl apply -f ./k8s/postgresql-secrets.yaml
	kubectl apply -f ./k8s/postgresql-cm.yaml
# other resources
	kubectl apply -f ./k8s

.PHONY: k8s-clean
k8s-clean:
	kubectl delete -f ./k8s
