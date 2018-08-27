default: stretch

all: stretch buster

update:
	@./update.sh

stretch:
	docker build --tag opxhub/gbp:stretch -f debian/stretch/base/Dockerfile .
	docker build --tag opxhub/gbp:stretch-dev debian/stretch

buster:
	docker build --tag opxhub/gbp:buster -f debian/buster/base/Dockerfile .
	docker build --tag opxhub/gbp:buster-dev debian/buster

.PHONY: default all update stretch buster
