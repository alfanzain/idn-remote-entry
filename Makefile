.PHONY: *

test:
	go test -v -timeout=1m30s ./...

run:
	-docker compose -f ./deploy/local/run/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/docker-compose.yml up --build

## command for Batha server which run old AMI Linux version. In there there is no "docker compose".
deploy-ec2:
	-make stop-ec2
	docker-compose -f ./deploy/aws/ec2/docker-compose.yml up --build -d

logs-ec2:
	docker-compose -f ./deploy/aws/ec2/docker-compose.yml logs -f

stop-ec2:
	docker-compose -f ./deploy/aws/ec2/docker-compose.yml down --remove-orphans	
