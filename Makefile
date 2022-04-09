.SILENT:

.PHONY: run
run:
	./scripts/run.sh

.PHONY: test
test:
	./scripts/test.sh

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: docker
docker:
	./scripts/docker.sh

.PHONY: stat
stat:
	./scripts/stat.sh
	