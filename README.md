# quick start
## create app

    serverless create -u https://github.com/gaku3601/serverless-go/ -p your_app_name

## deploy

    GOOS=linux go build -o bin/create ./func/create
    sls deploy

## stage切り替え

    sls deploy --stage prod

## stack削除

    serverless remove -v

# curlコマンド

    curl -X POST -H 'Content-Type:application/json' -d '{"title":"val"}' <URL>
