# Gorello

Gorrello is back-end for task management application, like Trello, written in go.  
It's deployed on [Heroku](https://friendly-drake-69422.herokuapp.com/).  
Swagger documentation is available under root (*/*) endpoint.

### Local deploy
For local deploy you need to have Docker installed.  
If you already have it, just run command below.
```bash
docker-compose up
```
If you experience permission problems try to run command with *sudo*.  
Application will be started on *localhost:8080*

### Testing
Tests also run in Docker. Use *test.sh* script to run all tests.