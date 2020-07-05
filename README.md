# monitoring-controller

Uses kubebuilder: https://book-v1.book.kubebuilder.io/

## Updating api/ types:

1. update `*_types.go` file 
1. `make manifests` - generate CRD YAML in `config/crd/bases/`
1. `make generate` - generate GO code

## Samples

- [HttpMonitoring](config/samples/monitoring.raisingthefloor.org_v1alpha1_httpmonitor.yaml)