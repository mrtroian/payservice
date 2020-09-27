# Payservice

### Building

Make sure you have Go compiler version 1.14.
Clone the repository, enter the folder and run command:

## Prerequesities

You need to have working SSL certificates that is registered in your system.
You can generate some with openssl:

	$ openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
	$ openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt

Note: app will export them automatically, make sure they are in configs/ssl/

## Running

To run the app use:

	$ make run

The `make run` will generate easyjson, build, export env and launch the app.
You can also `make` and launch the app manualy, don't forget to
set 'PAYSERVICE_CONFIGS_DIR' enviromental variable to configs/ folder:

	$ export PAYSERVICE_CONFIGS_DIR=$(pwd)/configs/
	$ ./bin/payservice

Note: There are only one config.yaml for the app, no dev or test configs regarded.
There are also only one mockconfig.yaml for mocker-server.

## Running tests

The `make test` will run unit tests, there're a few of them.
The `make test-integration` will run integration tests. Make sure to start
the app and the mocker-server before running integration tests.
Note: There is only one basic integration test, it only validates the response
to be 200. No data validation.

The flow:

	$ make run
	$ make run-mocker
	$ make integration-test

## Manual testing

After `make run` and `make run-mocker`, you can enjoy the API via Postman or whatever you like!
