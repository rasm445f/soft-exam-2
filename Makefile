up:
	docker-compose -f services/customerService/docker-compose.yml up -d
	docker-compose -f services/orderService/docker-compose.yml up -d
	docker-compose -f services/restaurantService/docker-compose.yml up -d
	docker-compose -f services/shoppingCartService/docker-compose.yml up -d
	docker-compose -f broker/docker-compose.yml up -d
	docker-compose -f apiGateway/docker-compose.yml up -d
	docker-compose -f monitoring/docker-compose.yml up -d

down:
	docker-compose -f services/customerService/docker-compose.yml down
	docker-compose -f services/orderService/docker-compose.yml down
	docker-compose -f services/restaurantService/docker-compose.yml down
	docker-compose -f services/shoppingCartService/docker-compose.yml down
	docker-compose -f broker/docker-compose.yml down
	docker-compose -f apiGateway/docker-compose.yml down
	docker-compose -f monitoring/docker-compose.yml down



run-customer:
	@echo "Running customerService..."
	$(MAKE) -C services/customerService run-go

run-order:
	@echo "Running orderService..."
	$(MAKE) -C services/orderService run-go

run-restaurant:
	@echo "Running restaurantService..."
	$(MAKE) -C services/restaurantService run-go

run-shoppingcart:
	@echo "Running shoppingCartService..."
	$(MAKE) -C services/shoppingCartService run-go

run-all:
	@trap 'echo "Stopping..."; kill 0' SIGINT SIGTERM; \
	( \
		$(MAKE) run-customer & \
		$(MAKE) run-order & \
		$(MAKE) run-restaurant & \
		$(MAKE) run-shoppingcart & \
		wait \
	)

docs-all:
	@echo "Generating docs..."
	$(MAKE) -C services/customerService docs
	$(MAKE) -C services/orderService docs
	$(MAKE) -C services/restaurantService docs
	$(MAKE) -C services/shoppingCartService docs

migrate-down-all:
	@echo "Migrating down..."
	$(MAKE) -C services/customerService migrate-down
	$(MAKE) -C services/orderService migrate-down
	$(MAKE) -C services/restaurantService migrate-down

migrate-up-all:
	@echo "Migrating up..."
	$(MAKE) -C services/customerService migrate-up
	$(MAKE) -C services/orderService migrate-up
	$(MAKE) -C services/restaurantService migrate-up

env-all:
	@echo "Creating .env files..."
	cp services/customerService/.env.example services/customerService/.env
	cp services/orderService/.env.example services/orderService/.env
	cp services/restaurantService/.env.example services/restaurantService/.env
	cp services/shoppingCartService/.env.example services/shoppingCartService/.env

