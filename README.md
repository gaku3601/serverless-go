# quick start

    serverless create -u https://github.com/gaku3601/serverless-go/ -p your_app_name

# deploy

    GOOS=linux go build -o bin/func ./func
    sls deploy

# stage切り替え

    sls deploy --stage prod

# stack削除

    serverless remove -v
