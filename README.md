# quick start
## create app

    serverless create -u https://github.com/gaku3601/serverless-go/ -p your_app_name

## deploy

    GOOS=linux go build -o bin/create ./func/create
    GOOS=linux go build -o bin/show ./func/show
    GOOS=linux go build -o bin/index ./func/index
    sls deploy

## stage切り替え

    sls deploy --stage prod

## stack削除

    serverless remove -v

# curlコマンド

    [create]
    curl -X POST -H 'Content-Type:application/json' -d '{"title":"val"}' <URL>
    [index]
    curl -X GET <URL>?start=1\&end=10
