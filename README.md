# GOODFOOD

Visualizer of clients and transactions

This repository contains code for the backend of GOODFOOD application. [In this link you can see code for web interface](https://github.com/rubbenpad/gofood)

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

POST `/data` Load data according date send

| Param |   Sample   |      Type      | Required |
| :---: | :--------: | :------------: | :------: |
| date  | 1603238400 | Unix timestamp |   true   |

### Buyers

GET `/buyers` All buyers in the platform

GET `/buyers/{{buyerId}}` Information about a buyer

## Technologies:

-   [Go](https://golang.org)
-   [Chi Router](https://github.com/go-chi/chi)
-   [DGraph](https://dgraph.io)
-   [Vue](https://vuejs.org) and [Vuetify](https://vuetifyjs.com)

## Quick start:

1. Clone this repository:

    `git clone git@github.com:rubbenpad/gofood.git`

2. Navigate into your new folder and run dgraph database

    > You need to have installed docker and docker-compose

    `cd gofood`

    `docker-compose up -d`

3. Navigate into src/ folder and start the application

    `cd src`

    `go run main.go`
