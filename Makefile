#!make
include .env
export $(shell sed 's/=.*//' .env)

test:
	env
build:
	@docker build -t ${APP_IMAGE_TAG} -f ./docker/app/Dockerfile .
	@docker image push ${APP_IMAGE_TAG}

push-image:
	@docker image push ${APP_IMAGE_TAG}

pull-image:
	@docker image pull ${APP_IMAGE_TAG}

stack-deploy:
	@docker stack deploy --compose-file=docker-stack.yml twitch-chat-archive

stack-down:
	@docker stack down twitch-chat-archive

app-service-logs:
	@docker service logs twitch-chat-archive_app