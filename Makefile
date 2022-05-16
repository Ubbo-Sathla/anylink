
all:  pull stop build start

pull:
	git pull origin self-signed-tls
	git reset --hard origin/self-signed-tls
stop:
	docker-compose down

pre_build:
	docker system prune -f

build: pre_build
	docker-compose build

pre_start:
	docker stop v2ray
start: pre_start
	docker-compose up -d