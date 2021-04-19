## build docker image

`docker build -t supervisord:v4.2.2 -f Dockerfile .`

## run supervisord

`docker run --name=supervisord --network=host --rm supervisord:v4.2.2`
