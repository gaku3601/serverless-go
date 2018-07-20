build:
	cd ./func && dep ensure	
	GOOS=linux go build -o bin/lambda ./func
deploy:
	@make build
	sls deploy

prod-deploy:
	@make build
	sls deploy --stage prod

remove:
	serverless remove -v
