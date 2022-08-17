<h2 align="center">API Go FAIR:</h2>
<p>
  <img alt="In Development" align="center" src="atWork.png" />
  <img alt="Version" src="https://img.shields.io/badge/version-1.00.0-blue.svg?cacheSeconds=2592000" />
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>

</p>

- Go API Fair. It also has a CRUD.

## Requirements/dependencies
- GO
- Docker
- Docker-compose
- MongoDB/SQL

## Getting Started

- After installing Go and setting up your GOPATH, 
- [How To install Go](https://github.com/larien/aprenda-go-com-testes/blob/master/primeiros-passos-com-go/instalacao-do-go.md) 


- [Clone project](https://github.com/marcovargas74/m74-fair-api.git)
```sh
git clone https://github.com/marcovargas74/m74-fair-api.git
```

- HOW TO RUN   
```sh
 cd mm74-fair-api

 ## start dockers
 make all

 ## stop dockers
 make stop
```

> :warning: **DB can take up to 3 minutes to start**: Be very careful here!


- Enter in project

```sh
cd m74-fair-api/src/fair
```

- Build e RUN golang project
```sh
 ## Run compiled project
	go run main.go
```

- Build golang project

```sh
go build -o main.go
```
- Run api(port 5000)
```sh
 ## Run compiled project
	go run main.go
```



## API Request

| Endpoint        | HTTP Method           | Description           |
| --------------- | :-------------------: | :-------------------: |
| `/status`       | `GET`                 | `Get status`          |
| `/all`          | `GET`                 | `Get All CPFs/CNPJs`  |
| `/cpfs/{cpf}`   | `GET`                 | `Check CPF`           |
| `/cpfs/{cpf}`   | `DELETE`              | `Delete CPF`          |
| `/cpfs`         | `GET`                 | `List All CPF`        |
| `/cnpjs/{cnpj}` | `GET`                 | `Check a CNPJ`        |
| `/cnpjs/{cnpj}` | `DELETE`              | `Delete CNPJ`         |
| `/cnpjs`        | `GET`                 | `List All CNPJ`       |


## Test endpoints API using curl

- #### Check status

`Request`
```bash
curl -i --request GET 'http://localhost:5000/status' \
```

`Response`
```json
{
    "num_total_query": 0,
    "up_time": 7.313309784066667,
    "start_time": "09-Jul-22 18:15:54"
}
```

- #### Listing ALL CPF and CNPJs

`Request`
```bash
curl -i --request GET 'http://localhost:5000/all'
```

`Response`
```json

[
    {
        "id":"bc78d0f0-4107-4b7e-a82e-d57a0167d2ca",
        "cpf":"682.511.941-99",
        "is_valid":true,
        "is_cpf":true,
        "is_cnpj":false,
        "created_at":"01-Jan-01 00:00:00"
    },

    {
        "id":"3b8416ad-cd38-471c-8259-e886c0aa91af",
        "cpf":"838.461.722-86",
        "is_valid":true,
        "is_cpf":true,
        "is_cnpj":false,
        "created_at":"01-Jan-01 00:00:00"
    }
    .
    .
    .
]

```

- #### Check and Creating new CPF Consult

`Request`
```bash
curl -i --request GET 'http://localhost:5000/cpfs/682.511.941-99' 
```

`Response`
```json
{
    "id": "28f0d8fa-f76f-47bd-bd65-58a3c4ee9c12",
    "cpf": "682.511.941-99",
    "is_valid": true,
    "is_cpf": true,
    "is_cnpj": false,
    "created_at":"09-Jul-22 18:50:56"
}
```
- #### Listing CPFs

`Request`
```bash
curl -i --request GET 'http://localhost:5000/cpfs'
```

`Response`
```json
[
    {
    "id":"5cf59c6c-0047-4b13-a118-65878313e329",
    "cpf":"111.111.111-11",
    "status":"isValid",
    "created_at":"2022-01-24T10:10:02Z"
    }
]
```

- #### Delete CPF Number

`Request`
```bash
curl -i --request DELETE 'http://localhost:5000/cpfs/682.511.941-99' 
```

- #### Check and Creating new CNPJ Consult

`Request`
```bash
curl -i --request POST 'http://localhost:5000/cnpj/73.212.132/0001-50' 
```

`Response`
```json
{
    "id":"7cf59c6c-0047-4b13-a118-65878313e329",
    "cnpj":"73.212.132/0001-50",
    "status":"isValid",
    "created_at":"2022-01-24T10:10:02Z"
}
```
- #### Listing CNPJs

`Request`
```bash
curl -i --request GET 'http://localhost:5000/cnpj'
```

`Response`
```json
[
    {
    "id":"7cf59c6c-0047-4b13-a118-65878313e329",
    "cnpj":"73.212.132/0001-50",
    "status":"isValid",
    "created_at":"2022-01-24T10:10:02Z"
    }
]
```

- #### Delete CNPJ Number

`Request`
```bash
curl -i --request DELETE 'http://localhost:5000/cnpj/73.212.132/0001-50' 
```


## Code status
- Development

## Next Steps
- Make a refactory
- Fix some bugs
- Add more tests
- make a interface web to test api

## Author
- Marco Antonio Vargas - [marcovargas74](https://github.com/marcovargas74)

## License
Copyright © 2022 [marcovargas74](https://github.com/marcovargas74).
This project is [MIT](LICENSE) licensed.
