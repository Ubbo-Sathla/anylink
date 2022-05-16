
all:  stop build start

stop:
	docker-compose down

pre_build:
	docker start v2ray
	docker system prune -f

build: pre_build
	docker-compose build

pre_start:
	docker stop v2ray
start: pre_start
	docker-compose up -d