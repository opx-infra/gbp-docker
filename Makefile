default: stretch

all: stretch

update:
	@./update.sh

stretch:
	docker build --tag opxhub/gbp:stretch -f debian/stretch/base/Dockerfile .
	docker build --tag opxhub/gbp:stretch-dev debian/stretch

.PHONY: default all update stretch
