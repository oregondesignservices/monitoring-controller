package conf

import (
	"errors"
	"github.com/oregondesignservices/monitoring-controller/internal/httpclient"
	"github.com/oregondesignservices/monitoring-controller/internal/metrics"
	"github.com/urfave/cli/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strings"
	"time"
)

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "metrics-addr",
			Value: ":8080",
			Usage: "The address the metric endpoint binds to.",
		},
		&cli.StringFlag{
			Name:  "namespace",
			Value: "monitoring",
			Usage: "which namespace the controller is running within",
		},
		&cli.DurationFlag{
			Name:  "http-client-timeout",
			Value: 29 * time.Second,
			Usage: "the http client timeout duration",
		},
		&cli.BoolFlag{
			Name:  "enable-leader-election",
			Usage: "Enable leader election for controller manager",
		},
		&cli.StringSliceFlag{
			Name:  "set-var",
			Usage: "set a global variable available to all requests. Format: 'key=value'",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose output",
		},
	}
)

type configuration struct {
	MetricsAddr          string
	Namespace            string
	HttpClientTimeout    time.Duration
	EnableLeaderElection bool
	GlobalRequestVars    map[string]string
}

func (c *configuration) UpdateFromCli(ctx *cli.Context) error {
	c.MetricsAddr = ctx.String("metrics-addr")
	c.Namespace = ctx.String("namespace")
	c.HttpClientTimeout = ctx.Duration("http-client-timeout")
	c.EnableLeaderElection = ctx.Bool("enable-leader-election")

	httpclient.Initialize(c.HttpClientTimeout)
	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	logger := ctrl.Log.WithName("configuration").WithName("UpdateFromCli")

	setvars := ctx.StringSlice("set-var")

	for _, v := range setvars {
		pieces := strings.SplitN(v, "=", 2)
		if len(pieces) != 2 {
			return errors.New("--set-var format must be 'key=value'")
		}
		c.GlobalRequestVars[pieces[0]] = pieces[1]
		logger.Info("added global request variable", "key", pieces[0], "value", pieces[1])
		metrics.GlobalVarsDetails.WithLabelValues(pieces...).Inc()
	}
	return nil
}

var GlobalConfig = &configuration{
	GlobalRequestVars: make(map[string]string),
}
