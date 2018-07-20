# quick start
## create app

    serverless create -u https://github.com/gaku3601/serverless-go/ -p your_app_name

## deploy

    make deploy

## deploy production

    make prod-deploy

# curlコマンド

    [create]
    curl -X POST -H 'Content-Type:application/json' -d '{"title":"val"}' <URL>
    [index]
    curl -X GET <URL>?start=1\&end=10
    [update]
    curl -X PATCH -H 'Content-Type:application/json' -d '{"title":"val4"}' <URL>/{id}
    [destroy]
    curl -X DELETE <URL>/{id}
    [show]
    curl -X GET <URL>/{id}
