#!make
include .env
export $(shell sed 's/=.*//' .env)

test:
	env
build:
	@docker build -t ${APP_IMAGE} -f ./docker/app/Dockerfile .

push-image:
	@docker image push ${APP_IMAGE}

pull-image:
	@docker image pull ${APP_IMAGE}

stack-deploy:
	@docker stack deploy --compose-file=docker-stack.yml twitch-chat-archive

stack-down:
	@docker stack down twitch-chat-archive

app-service-logs:
	@docker service logs twitch-chat-archive_app