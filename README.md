# go-template

Simple binary to use Go's [template engine](https://pkg.go.dev/text/template).

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/airfocusio/go-template:latest

export NAME=tom
echo "Hello, {{ .Env.NAME }}" | docker run --rm -i ghcr.io/airfocusio/go-template:latest
```
