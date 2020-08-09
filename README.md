# monitoring-controller

Install CRDs into kubernetes that periodically monitor resources. Currently handles HTTP, but
could be extended to handle more.

Data is exported as prometheus metrics, which may be used to send alerts.

## CustomResourceDefinitions

- [HttpMonitor](config/crd/bases/monitoring.raisingthefloor.org_httpmonitors.yaml)

## Examples

See [samples](config/samples).

## Available Metrics

See [metrics.go](internal/metrics/metrics.go).

## Grafana Dashboard

The grafana dashboard may be found in the kustomize-based [deployment repo](https://github.com/oregondesignservices/deploy-monitoring-controller/blob/master/resources/grafana/main-dashboard.json).

## Kubebuilder

This project uses [kubebuilder](https://book.kubebuilder.io/) as a way of generating code.

To install kubebuilder:
```shell script
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch}/bin/kubebuilder /usr/local/
```

### Updating Existing Monitoring CRDs:

1. update `*_types.go` file 
1. `make manifests` - generate CRD YAML in `config/crd/bases/`
1. `make generate` - generate GO code

### Adding New Monitoring CRDs

1. `kubebuilder create api --group monitoring.raisingthefloor.org --version v1alpha1 --kind [new kind]`
1. fill in `api/v1alpha1/[new kind]_types.go`
1. implement `controllers/[new kind]_controller.go`
1. `make manifests generate`

## Development

This project uses [KIND](https://kind.sigs.k8s.io/) and [tilt](https://tilt.dev/) for rapid development.

For setup:
```shell script
export KUBECONFIG=kind-kubeconfig.yaml

# Sets up a KIND cluster and a local registry connected to it
sh kind-up.sh

# start tilt. Hit space after to open tilt in your browser (recommended)
# Local changes will automatically be built and pushed to kube.
tilt up

# To shut down:
tilt down
sh kind-down.sh
```