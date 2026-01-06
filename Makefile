.PHONY: all clean default run build update check pc test

default: check

run:
	SERVICE_DSN=test go run ./cmd/template-go

build:
	go build -o ./dist/ ./...

goreleaser:
	goreleaser --clean --snapshot --skip=publish

update:
	go get -u -t ./...
	go mod tidy
	go mod verify
	prek auto-update
	pinact run -update

check: pc test
pc:
	prek run -a
test:
	TESTCONTAINERS=1 go test -v -race -covermode=atomic -coverprofile=coverage.txt ./...

bumped:
	git cliff --bumped-version

# make release TAG=$(git cliff --bumped-version)-alpha.0
release: check
	git cliff -o CHANGELOG.md --tag $(TAG)
	prek run --files CHANGELOG.md || prek run --files CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "chore(release): prepare for $(TAG)"
	git push
	git tag -a $(TAG) -m "chore(release): $(TAG)"
	git push origin $(TAG)
