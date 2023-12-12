# Context Parameter

## Syntax



## When to use context parameters

Always use `context.Context` for functions that interact with the system
external.

* update cache value di redis ? use `context.Context`
* query ke postgres / mongodb ? use `context.Context`
* call API dari 3rd party ? use `context.Context`
* dial gRPC host / service ? use `context.Context`
* etc

> context. For. All. external. call.

## Hati-hati menggunakan hardcoded timeout context

Consider the following example:

```go
func UpdateSomething(ctx context.Context, param string) error {
  // add span to ctx...
  span, ctx := apm.StartSpan(ctx, "myspan", spanType)
  defer span.End()

  rCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
  defer cancel()

  // continue function...
}
```

`context.WithTimeout` will make a copy of the existing context, with a certain timeout.
The above line means that for 1 function block, the timeout is 10 seconds.

This becomes a problem if the timeout is created in each function. Example:
* `controller.UpdatePhone -> service.UpdatePhone -> mongo.UpdatePhone`
* If each function has a timeout of 10s, it means that the total request is possible to run for 30s.
* Even though the timeout from the user / android is a maximum of 10s
* There is a possibility on the user's cellphone timeout, but on the server the program is still being executed.

Ideally the timeout is set from the echo / grpc server middleware side, so 1x is enough. Furthermore, functions that require context just take the parameters of the previous function.

If there is already a `context.Context` parameter of the function, there is no need to create a new context with a different timeout.

## Reference

* https://golang.org/pkg/context/#WithTimeout
* [golang context guide]( https://golangbyexample.com/using-context-in-golang-complete-guide/ )
