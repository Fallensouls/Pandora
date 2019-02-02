module github.com/Fallensouls/Pandora

replace (
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => github.com/golang/crypto v0.0.0-20181127143415-eb0de9b17e85
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/sse v0.0.0-20190125020943-a7658810eb74 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-ini/ini v1.41.0
	github.com/go-xorm/core v0.6.0
	github.com/go-xorm/xorm v0.7.1
	github.com/golang/protobuf v1.2.0 // indirect; indirectgo
	github.com/lib/pq v1.0.0
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/ugorji/go/codec v0.0.0-20190126102652-8fd0f8d918c8 // indirect
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
