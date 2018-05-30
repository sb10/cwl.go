default: test

# go get -u github.com/stretchr/testify/assert
test:
	@go test --count 1

# go get -u gopkg.in/alecthomas/gometalinter.v2
# gometalinter.v2 --install
lint:
	@gometalinter.v2 --vendor --aggregate --deadline=120s ./... | sort

lintextra:
	@gometalinter.v2 --vendor --aggregate --deadline=120s --disable-all --enable=gocyclo --enable=dupl ./... | sort

.PHONY: test lint lintextra
