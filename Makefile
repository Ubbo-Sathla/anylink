
all:  pull stop build start

pull:
	git pull origin self-signed-tls -f
	git reset --hard origin/self-signed-tls
stop:
	docker-compose down

pre_build:
	docker system prune -f

build: pre_build
	docker-compose build

start:
	docker-compose up -d

restart: stop start

self:
	sed -i 's/bonc_cert.crt/vpn_cert.crt/g' /opt/anylink-conf/server.toml
	sed -i 's/bonc_cert.key/vpn_cert.key/g' /opt/anylink-conf/server.toml

bonc:
	sed -i 's/vpn_cert.crt/bonc_cert.crt/g' /opt/anylink-conf/server.toml
	sed -i 's/vpn_cert.key/bonc_cert.key/g' /opt/anylink-conf/server.toml