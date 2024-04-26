

all: server client models


server:
	oapi-codegen --config ./server.yaml ./b7s-swagger.yaml

client:
	oapi-codegen --config ./client.yaml ./b7s-swagger.yaml

models:
	oapi-codegen --config ./models.yaml ./b7s-swagger.yaml

