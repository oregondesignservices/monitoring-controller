# monitoring-controller

Uses kubebuilder: https://book-v1.book.kubebuilder.io/

## Updating api/ types:

1. update `*_types.go` file 
1. `make manifests` - generate CRD YAML in `config/crd/bases/`