# Echo HTTP API Design

There are many ways to implement http json web server in go. Can use std
lib `net/http`, or go web-framework (iris, fiber, gin, gorilla or echo).

Here will be explained related to the standard use of echo webframework.

## Tasks from http JSON layer / controller / handler

* Expose the function of `internal/service`. 1 function endpoint handler can
  mapped to 1 service method.
* Validate input from the user. The validation did not check
  external db/storage level. So just check the payload format. For validation
  those who check external API, db or other storage, can be placed at level
  http middleware & `internal/service` itself.
* Build parameters for the service method to be called.
* Call service method exposed
* Build response is returned to the user according to the results of the service method.
* Handle returned error (if any).