module github.com/project-eria/go-wot

go 1.19

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gofiber/fiber/v2 v2.38.1
	github.com/gofiber/websocket/v2 v2.1.0
	github.com/rs/zerolog v1.28.0
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/fasthttp/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.40.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.2.1
	v0.2.0
	v0.1.4
	v0.1.3
	v0.1.2
	v0.1.1
	v0.1.0
	v0.0.0
)
