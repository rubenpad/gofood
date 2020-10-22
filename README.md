# GOODFOOD

Visualizer of clients and transactions

This repository contains code for the backend of GOODFOOD application. [In this link you can see code for web interface](https://github.com/rubbenpad/vuefood)

## Requirements:

1. Get data from external sources and format it to store in a graph-based database

    - Products
    - Clients
    - Transactions

2. List all clients of the platform

3. Get information about a specific client by his ID

    - Transactions history
    - Other clients using the same IP
    - Recommendations about products that other people also buy

4. Create a web interface to visualize the data

    - Select a date to load data from
    - List all clients
    - See information about specific client

## API

### Load data

**POST** `/data` Load data according date sent in query params

| Param |   Sample   |      Type      | Required |
| :---: | :--------: | :------------: | :------: |
| date  | 1603238400 | Unix timestamp |   true   |

Example:

`http://localhost:3000/data?date=1603238400`

### Buyers

**GET** `/buyers` All buyers in the platform

Example:

`http://localhost:3000/buyers`

**GET** `/buyers/{{buyer_id}}` Information about a buyer

|  Param   |  Sample  |  Type  | Required |
| :------: | :------: | :----: | :------: |
| buyer_id | 6d910e7c | string |   true   |

Example:

`http://localhost:3000/buyers/6d910e7c`

## Technologies:

-   [Go](https://golang.org)
-   [Chi Router](https://github.com/go-chi/chi)
-   [DGraph](https://dgraph.io)
-   [Vue](https://vuejs.org) and [Vuetify](https://vuetifyjs.com)

## Quick start:

1. Clone this repository:

    `git clone git@github.com:rubbenpad/gofood.git`

2. Navigate into your new folder `gofood/` and create a `.env` file with the same fields that are in the `.env.example` file:

    `cd gofood`

    ```sh
        # External endpoint to request data
        BASE_URL=
        # Database host localhost:9080 see dgraph docs
        DGRAPH_HOST=
        # YES the first time you run the application to sync
        # the schema in your database
        SETUP_DB=
    ```

3. Run DGraph database:

    > You need to have installed docker and docker-compose

    `docker-compose up -d`

4. Navigate into `src/` folder compile and start the application

    `cd src`

    `go run main.go`
