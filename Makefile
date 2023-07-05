run:
	@docker-compose -f ./build/docker-compose.yml build
	@docker-compose -f ./build/docker-compose.yml up -d

up:
	@docker-compose -f ./build/docker-compose.yml up -d

down:
	@docker-compose -f ./build/docker-compose.yml down

destroy:
	@docker-compose -f ./build/docker-compose.yml down -v

log:
	@docker-compose --project-directory=./build logs