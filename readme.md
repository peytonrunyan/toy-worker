This is a toy example showing a worker that can either write to or consume from a `repository`such as RabbitMQ, Reddis, or MongoDB. It currently only supports RabbitMQ.

The worker will print messages from the repository when "consuming" and will write a simple message 

### To run
1. pull the RabbitMQ docker image `docker pull rabbitmq:3-management`
2. run the docker image `docker run --rm -it --hostname my-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management`
3. run `./dataprocessor --produce` in one terminal window to begin running the toy producer
4. run `./dataprocessor` in another terminal window to begin running the toy consumer