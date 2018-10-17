.PHONY: run keys test docker

docker: keys
	docker build -t gochallenge .

gochallenge: *.go
	go build

run: *.go
	go run *.go

.venv:
	python3.7 -m venv .venv
	.venv/bin/pip install requests pytest mock

test: .venv
	.venv/bin/pytest tests.py -v

keys:
	mkdir -p etc
	openssl genrsa 1024 > etc/key.priv
	cat etc/key.priv | openssl rsa -pubout > etc/key.pub
