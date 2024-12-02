up:
	docker-compose -f services/customerService/docker-compose.yml up -d
	docker-compose -f services/orderService/docker-compose.yml up -d
	docker-compose -f services/restaurantService/docker-compose.yml up -d
	docker-compose -f services/shoppingCartService/docker-compose.yml up -d
	docker-compose -f broker/docker-compose.yml up -d

down:
	docker-compose -f services/customerService/docker-compose.yml down
	docker-compose -f services/orderService/docker-compose.yml down
	docker-compose -f services/restaurantService/docker-compose.yml down
	docker-compose -f services/shoppingCartService/docker-compose.yml down
	docker-compose -f broker/docker-compose.yml down



run-customer:
	@echo "Running customerService..."
	$(MAKE) -C services/customerService run

run-order:
	@echo "Running orderService..."
	$(MAKE) -C services/orderService run

run-restaurant:
	@echo "Running restaurantService..."
	$(MAKE) -C services/restaurantService run

run-shoppingcart:
	@echo "Running shoppingCartService..."
	$(MAKE) -C services/shoppingCartService run

run-all:
	@trap "kill 0" SIGINT; \
	$(MAKE) run-customer & \
	$(MAKE) run-order & \
	$(MAKE) run-restaurant & \
	$(MAKE) run-shoppingcart & \
	wait

docs-all:
	@echo "Generating docs..."
	go-swagger-merger -o apiGateway/swagger-config.json  services/customerService/docs/swagger.json services/orderService/docs/swagger.json services/restaurantService/docs/swagger.json services/shoppingCartService/docs/swagger.json