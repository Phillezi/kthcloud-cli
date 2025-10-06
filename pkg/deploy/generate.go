package deploy

//go:generate go run gurl.go https://api.cloud.cbh.kth.se/deploy/v2/docs/doc.json openapi.json /v2/snapshots /v2/vmActions

//go:generate go tool oapi-codegen --config=oapi.config.yaml openapi.json
