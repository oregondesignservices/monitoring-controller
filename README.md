# monitoring-controller

Install CRDs into kubernetes that periodically monitor resources

## Development

Uses kubebuilder: https://book-v1.book.kubebuilder.io/

```shell script
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch}/bin/kubebuilder /usr/local/
```

## Updating api/ types:

1. update `*_types.go` file 
1. `make manifests` - generate CRD YAML in `config/crd/bases/`
1. `make generate` - generate GO code

## Samples

- [HttpMonitoring](config/samples/sample2.yaml)