LOGFILE=$(LOGPATH) `date +'%A-%b-%d-%Y-%H-%M-%S'`

.PHONY: hp
hp: ## 🌱 This help.💙
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: cm
cm: ## 🌱 git commit 💙
	@echo '************👇  run command 👇************'
	git add .
	git commit -m "🌱sample-blockchain-grpc💙-${LOGFILE}"
	git push -u origin main


.PHONY: server
server: ## 🌱 run application 💙
	go run  server/main.go

.PHONY: list
list: ## 🌱 view list all block 💙
	go run  client/main.go

.PHONY: add
add: ## 🌱 create block in blockchian 💙
	go run  client/main.go