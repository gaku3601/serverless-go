build:
	cd ./func/create && dep ensure	
	cd ./func/index && dep ensure	
	cd ./func/update && dep ensure	
	GOOS=linux go build -o bin/create ./func/create
	GOOS=linux go build -o bin/index ./func/index
	GOOS=linux go build -o bin/update ./func/update
deploy:
	@make build
	sls deploy

prod-deploy:
	@make build
	sls deploy --stage prod

remove:
	serverless remove -v
