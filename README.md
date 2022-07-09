# ⚙️ migration-customers

[![Build](https://github.com/diogoalbuquerque/migration-customers/actions/workflows/build.yml/badge.svg)](https://github.com/diogoalbuquerque/migration-customers/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/diogoalbuquerque/migration-customers/branch/main/graph/badge.svg?token=yvpL6uN67X)](https://codecov.io/gh/diogoalbuquerque/migration-customers)

## Project

This project consists of a processing function for creating a list of released clients in a new database,
having as source the customers already registered in the legacy database.

The flow starts with a scheduled execution through an AWS event.

When starting, eligible customers are founded.

When finding these customers:

- Each customer found will be verified if it is a person or legal entity;
- Customers validated in the previous step will be added into the new database;
- Customers added in the previous item will have their processing time scheduled so they are not eligible for
  new processing.

### Usage

Requires an event to start.

## Technologies used

This function was built using the following technologies:

- [aws-lambda-go - v1.32.1](https://github.com/aws/aws-lambda-go);
- [aws-sdk-go - v1.44.51](https://github.com/aws/aws-sdk-go);
- [aws-xray-sdk-go - v1.7.0](https://github.com/aws/aws-xray-sdk-go);
- [cleanenv v1.3.0](https://github.com/ilyakaznacheev/cleanenv);
- [go_ibm_db - v0.4.1](https://github.com/ibmdb/go_ibm_db);
- [go-sqlmock - v1.5.0](https://github.com/DATA-DOG/go-sqlmock);
- [mysql - v1.6.0](https://github.com/go-sql-driver/mysql);
- [testify v1.7.5](https://github.com/stretchr/testify);
- [zerolog v1.27.0](https://github.com/rs/zerolog).

## Technical operation

### How to install

#### Configuring DB2:

All instructions can be found here: [go_ibm_db](https://github.com/ibmdb/go_ibm_db)

But I'll put the step by step the way I understand to make it easier:

##### Step 01:

Download the dependency to your library:

```sh
go get -d github.com/ibmdb/go_ibm_db
```

At the end you will see that it has been installed in:

**{GO_PATH}/src/github.com/ibmdb/go_ibm_db**

##### Step 02:

Install ODBCCLI, this is a simple process, the library itself created a file for you to do this
automatically.

Look for:

**${GO_PATH}/src/github.com/ibmdb/go_ibm_db/installler**

and run:

```sh
go run setup.go
```

##### Step 02:

After CLIDRIVER is downloaded, you need to set the following environment variables:

| Variable | Value |
| ------ | ------ |
| DB2HOME | ${GO_PATH}/src/github.com/ibmdb/go_ibm_db/installer/clidriver |
| CGO_CFLAGS | I$DB2HOME/include |
| CGO_LDFLAGS | L$DB2HOME/lib |

❗ LINUX

| Variable | Value |
| ------ | ------ |
| LD_LIBRARY_PATH | ${GO_PATH}/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib |

❗ MAC

IF YOU ARE USING A recent OSX, this variable must be set at runtime

| Variable | Value |
| ------ | ------ |
| DYLD_LIBRARY_PATH | ${GO_PATH}/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib |

##### Step 03:

Add the license file in the CLIDRIVER folder.

Look for:

**${GO_PATH}/src/github.com/ibmdb/go_ibm_db/installer/clidriver/license**

and add the DB2 license file (usually the name is: **db2consv_ee.lic**)

## Main variables to be defined

| Variable           | Variable description                               | Default value      |
|--------------------|----------------------------------------------------|--------------------|
| LOG_LEVEL          | Log level.                                         | info               |
| MYSQL_OPEN_CONN_MAX | Maximum number of connections that can be opened.  | 100                |
| MYSQL_IDLE_CONN_MAX | Maximum number of connections that can be stopped. | 2                  |
| MYSQL_LIFE_CONN_MAX | Maximum connection time in seconds.                | 10                 |
| DB2_OPEN_CONN_MAX  | Maximum number of connections that can be opened.  | 100                |
| DB2_IDLE_CONN_MAX  | Maximum number of connections that can be stopped. | 2                  |
| DB2_LIFE_CONN_MAX  | Maximum connection time in seconds.                | 10                 |
| AWS_REGION_NAME    | Region where the secrets manager is located.       | DEFINE_REGION_NAME |
| SECRETS_MANAGER    | Name of the secret manager used in the application | SECRETS_MANAGER    |
| DATASOURCE_LIMIT   | Number of records that will be loaded per event.   | 100                |

#### Project build:

To upload a function you need to create and build it like linux operating.
So use the following command in the project folder:

```sh
$ GOOS=linux GOARCH=amd64 go build -o MIGRATION_CUSTOMERS ./
```

or you can use Makefile.

```sh
$ make build
```

❗ Attention️

- After the build, you need to create the .zip of the generated file and then upload this .zip file;
- The Handler is the name of the generated .zip file.

## Comments

- This project is a part of a proof of concept used in the company, so parts of them were omitted as well as domain and
  data model information.
- As the execution was sporadic, it was chosen to use it as a lambda, so there was no need to allocate a container and
  the execution time was acceptable for the solution.
- Another reason chosen was the possibility for the operator to request the execution whenever necessary.