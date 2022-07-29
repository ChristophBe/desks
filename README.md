# Desks
In current times the hybrid working gets a common practice. This leads to the need to manage who will need a desks, that enough free desks are available all times. This porject has the goal to build a desk booking system to comforably manage the available desks.

## Usage
Start the application use `docker-compose --env-file example.env up` add rooms to the database and start booking desks. The web app will be available on `http://localhost:5000`.

## Manualy run the Project
### Preconditions
- Install [go](https://go.dev/dl/) to be able to build and run the api.
- Install [yarn](https://classic.yarnpkg.com/lang/en/docs/install/) to build and run the web api.
- Create a environment file based on the `example.env` file. For Authentication you need  tp add the urls and credetial for an openid-connect provider. You can use for example [keycloak](https://www.keycloak.org/) as openid-connect provider or use a identity service of a cloud provider.

### Database
```bash
docker-compose --env-file example.env db
```

### Api
```bash
source example.env
cd desks-api/
go run main.go
```

### Web-App
```bash
cd desks-web-app/
yarn serve
```


