# go-template

Simple binary to use Go's [template engine](https://pkg.go.dev/text/template). Just pipe the template via stdin and retrieve the rendered output from stdout.

## Values

| Key                              | Value                                          |
|----------------------------------|------------------------------------------------|
| `.Now`                           | Current timestamp                              |
| `.Val.some`                      | Retrieve any value provided via CLI argument   |
| `.Env.SOME`                      | Retrieve any environment variable              |
