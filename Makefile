ifneq (,$(wildcard ./.env))
  include .env
  export
endif

BUILD = $(DB)
BIN = $(HOME)/bin
SCHEMA = schema.sql

CONFIG_DIR=$(HOME)/.config/msgapi
CONFIG_FILE=$(CONFIG_DIR)/msgapirc

build:
	go build

install: build
	mv $(BUILD) $(BIN)

clean:
	$(RM) ./$(BUILD)

config:
	test -d $(CONFIG_DIR) || mkdir -p $(CONFIG_DIR)
	test -f $(CONFIG_FILE) || touch $(CONFIG_FILE)
	chmod 600 $(CONFIG_FILE)

db:
	sudo -u postgres createdb -O $(USER) $(DB)
	sudo -u postgres psql -c 'CREATE EXTENSION citext;' $(DB)
	psql -U $(USER) -f $(SCHEMA) $(DB)

csv:
	psql -c "\COPY msgs TO STDOUT CSV HEADER" $(DB)

.PHONY: build db install
