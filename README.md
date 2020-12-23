# Replicated Log

## Pre-requirements
- installed [docker](https://docs.docker.com/engine/install/)
- installed [docker-compose](https://docs.docker.com/compose/install/)

## How to run
- in the project root run - `docker-compose up`
- master node link: [http://localhost:7085](http://localhost:7085)
- secondary 1 node link: [http://localhost:7086](http://localhost:7086)
- secondary 2 node link: [http://localhost:7087](http://localhost:7087)

## How to test
- to add a message: make `POST` request to [http://localhost:7085/messages](http://localhost:7085/messages) 
  with header `Content-Type: application/json` and body
  ```
  {
    "body": "Message text",
    "w": 1
  }
  ```
  where `body` is string and contains message text and `w` is integer 
  represents write concern
- to get a list of messages: make `GET` request to one of available 
  nodes (master or secondary) to the `/messages` url
- see requests examples in the `api-requests.http` file

## Contributors
- Ihor Semeniuk
- Pavlo Sliusar
- Dmytro Ilchuk
