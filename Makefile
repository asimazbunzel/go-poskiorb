build:
	go build -o bin/orbits ./cmd/go-orbits

test:
	bin/orbits -C test/config.yaml
