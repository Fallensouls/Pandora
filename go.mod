module github.com/go-pandora/core

replace (
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => github.com/golang/crypto v0.0.0-20181127143415-eb0de9b17e85
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
)

require (
	github.com/Fallensouls/Pandora v0.0.0-20190312103849-598afe9fa638 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.3.0
	github.com/go-pandora/pkg v0.0.0-20190313091716-21e39597bac5
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/go-xorm/core v0.6.0
	github.com/go-xorm/xorm v0.7.1
	github.com/lib/pq v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85
	gopkg.in/yaml.v2 v2.2.2
)
