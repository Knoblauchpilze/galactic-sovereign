
COLOR_HIGHLIGHT_BLUE=\e[1;36m
COLOR_HIGHLIGHT_GREEN=\e[1;32m
COLOR_CLEAR=\e[0m

all: setup install

install:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Building apps...$(COLOR_CLEAR)"
	cd cmd/server && make install
	@echo "$(COLOR_HIGHLIGHT_GREEN)Success!$(COLOR_CLEAR)"

setup:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Setup project...$(COLOR_CLEAR)"
	@./scripts/setup.sh
	@echo "$(COLOR_HIGHLIGHT_GREEN)Success!$(COLOR_CLEAR)"

clean:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Cleaning apps...$(COLOR_CLEAR)"
	@cd bin && ls | grep -v .gitignore | xargs rm -fr
	@cd cmd/server && make clean
	@echo "$(COLOR_HIGHLIGHT_GREEN)Success!$(COLOR_CLEAR)"

# https://stackoverflow.com/questions/3931741/why-does-make-think-the-target-is-up-to-date
.PHONY: test
test:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Testing internal...$(COLOR_CLEAR)"
	./scripts/test.sh --dir internal
	@echo "$(COLOR_HIGHLIGHT_BLUE)Testing pkg...$(COLOR_CLEAR)"
	./scripts/test.sh --dir pkg
	@echo "$(COLOR_HIGHLIGHT_GREEN)Success!$(COLOR_CLEAR)"

# https://stackoverflow.com/questions/6273608/how-to-pass-argument-to-makefile-from-command-line
run:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Running $(filter-out $@,$(MAKECMDGOALS))...$(COLOR_CLEAR)"
	@cd bin && ./run.sh --app-name $(filter-out $@,$(MAKECMDGOALS))

docker:
	@echo "$(COLOR_HIGHLIGHT_BLUE)Building docker...$(COLOR_CLEAR)"
	docker build --tag user-service .

%:
	@:
