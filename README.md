# go-template

Simple binary to use Go's [template engine](https://pkg.go.dev/text/template). Just pipe the template via stdin and retrieve the rendered output from stdout.

## Values

| Key                              | Value                                                       |
|----------------------------------|-------------------------------------------------------------|
| `.Val.some`                      | Retrieve any value provided via CLI argument or yaml file   |
| `.Env.SOME`                      | Retrieve any environment variable                           |
