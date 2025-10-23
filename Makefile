setup:
	@echo "--- Setup and generated config yaml files ---"
	@mkdir -p config/resources
	@cp -r config/example/*.yml config/resources/

api:
	@echo "--- running api server in dev mode ---"
	@go run main.go api

migrate-up:
	@echo "--- running db migration ---"
	@go run main.go db migrate up $(step) --schema=$(schema)

migrate-down:
	@echo "--- running db migration ---"
	@go run main.go db migrate down $(step) --schema=$(schema)

migrate-create:
	@echo "--- creating db migration ---"
	@migrate create -ext sql -dir db/migrations/postgres -seq $(name)

build: setup
	@echo "--- Building binary file ---"
	@go build -o ./main main.go

download: 
	@echo "Downloading .tar file from joindiskon"
	@sudo rsync -avz root@62.72.59.72:/root/dplk/ /Users/raymondsugiarto/Documents/project/dplk/dplk-be/server/
	@sudo rclone sync /Users/raymondsugiarto/Documents/project/dplk/dplk-be/server/ dplk:/home/adminnss/dplk/
# 	@sudo rclone sync /Users/raymondsugiarto/Documents/project/dplk/dplk-be/server/ dplk-prod-26:/home/adminnss/dplk/
# 	@sudo rsync -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" -avz root@62.72.59.72:/root/dplk/ /Users/raymond.sugiarto/Documents/mine/p/dplk-be/server/
# 	@sudo rclone sync /Users/raymond.sugiarto/Documents/mine/p/dplk-be/server/ dplk:/home/adminnss/dplk/
