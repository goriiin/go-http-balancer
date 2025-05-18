

all_up:
	sudo docker compose up -d

rebuild_backend:
	sudo docker stop backend_group1_1 && sudo docker rm backend_group1_1 && sudo docker rmi  go-http-balancer-backend_group1_1 && sudo docker compose up -d

rebuild_balancer:
	sudo docker stop http_balancer && sudo docker rm http_balancer && sudo docker rmi go-http-balancer-balancer && sudo docker compose up -d

rebuild_all: rebuild_backend  rebuild_balancer

logs_back:
	sudo docker logs backend_group1_1

logs_balancer:
	sudo docker logs http_balancer

clean:
	go clean

.PHONY: clean