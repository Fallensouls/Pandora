# Pandora
Pandora is a simple web service API implemented with golang.

## Installation
`go get github.com/Fallensouls/Pandora`

## How to use
### Required
* PostgreSQL
* Redis
### Configuration
You can use your own configuration to run Pandora. For example:
```
service:
  name: Pandora

server:
  run_mode: debug  # debug, release or test
  port: 8080
  read_timeout: 60  # 60s
  write_timeout: 60

database:
  type: postgres
  name: postgres
  user: postgres
  password: *******
  host: 127.0.0.1
  port: 5432

redis:
  host: 127.0.0.1
  port: 6379
  password: *******
  
jwt:
  signing_algorithm: HS256  # HS256, HS384 or HS512
  secret: *******
  timeout: 60               # 60min
  issuer: Fallensouls
``` 
## Features
- [x] Restful API
- [x] JWT-based authentication
- [x] Yaml Configuration
- [ ] OAuth
- [ ] Swagger
- [ ] Log
- [ ] Docker
- [ ] Pandora-pkg
    - [ ] CAPTCHA
    - [ ] Email
    - [ ] SMS
    - [ ] QR Code
    
## Packages we use
* HTTP Router   [gin](https://gin-gonic.github.io/gin/) - [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* ORM   [xorm](http://xorm.io) - [github.com/go-xorm/xorm](https://github.com/go-xorm/xorm)
* Redis [github.com/go-redis/redis](https://github.com/go-redis/redis)
* YAML  [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)
* JWT   [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)

## Author
[Fallensouls](https://twitter.com/lu_tju?s=09) - I really love golang which changes my code style and thinking in programming.
Hope that everyone could enjoy golang!