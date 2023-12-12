# Production Service Checklist

## Here are the expectations for each go-service before deploying to production

- [ ] No error or warning from `golangci-lint`, template config
  use [template golangci.yml config](./golangci.template.yml)
- [ ] Can be built
- [ ] Error output does not stutter, check [error handling](./1_2_error_handling.md) &
  [logging](./5_logging.md). So that log messages to the point & not wasteful of storage.
- [ ] Implement context parameter sharing for each external call, check [context parameter](./1_4_context_parameter.md)
- [ ] Implement healthcheck , can refer to [http-restapi](./2_echo_api_design.md) &
  [gRPC-service](./3_grpc_api_design.md)
- [ ] Graceful shutdown. Every time before the app is killed, give a chance for each request
  / an ongoing process to complete within a specified timeout.
  In addition, aga can run optimally on a preemptible VM (cheaper).
  The logic is roughly as follows:
    * API/service is handling *N*-requests
    * API/service suddenly receive kill signal (*Ctrl-C, sigkill, etc*)
    * API/service return healthcheck NOK, then wait for *t*-seconds
      so that ongoing requests can be completed.
    * API/service exit program

## Reference

* [GolangCI-Lint linter](https://golangci-lint.run)
* [go worker process cancellation](https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/)
* [http graceful shutdown](https://www.rodrigoaraujo.me/posts/golang-pattern-graceful-shutdown-of-concurrent-events/)
* [GCP Preemptible VM](https://cloud.google.com/compute/docs/instances/preemptible)