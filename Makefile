DIST ?= stretch
DEBIANS := jessie stretch buster sid

ifneq ($(filter $(DIST),$(DEBIANS)),)
DIR := debian
else
DIR := ubuntu
endif

default: update
	docker build --tag opxhub/gbp:$(DIST) -f $(DIR)/$(DIST)/base/Dockerfile .
	docker build --tag opxhub/gbp:$(DIST)-dev $(DIR)/$(DIST)

update: Dockerfile-debian-base.template Dockerfile-debian-dev.template
	./update.sh

.PHONY: default update
