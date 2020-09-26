# Payservice

### About the bugs:
* Configuration exports not well tested, may break in some cases
* Config is not validated (yet)
* Be careful when running with your SSL certificates via `make run`

### Building

Make sure you have Go compiler version 1.14.
Clone the repository, enter the folder and run command:

## Prerequesities

You need to have working SSL certificates that is registered in your system.
You can generate some with openssl.
Before launching the app, please set SSL_KEY and SSL_CERT env variables.

Note: Makefile can export them automatically, make sure they are in configs/ssl/localhost.*

	$ export SSL_KEY=path/to/certificate.key
	$ export SSL_CERT=path/to/certificate.crt

## Running

To run the app use:

	$ make run

The `make run` will generate easyjson, build, export config and launch the app.
You can also `make` and launch the app manually by passing `config.yaml`:

	$ make run
	$ ./bin/payservice -config path/to/config.yaml

## Running tests

The `make test` will run unit tests, there're a few of them.
The `make test-integration` will run integration tests. Make sure to start
the app and the mocking server before running integration tests.

The flow:

	$ make run
	$ make mocker-run
	$ make integration-test

## Manual testing

After `make run` and `make mocker-run`, you can enjoy the API via Postman or whatever you like!