package cmd

import (
	"io"
	"strings"

	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"

	"github.com/jenkins-x/jx/pkg/jx/cmd/templates"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/util"
)

const (
	defaultSentryReleaseName = "sentry"
)

var (
	createAddonSentryLong = templates.LongDesc(`
		Creates the sentry addon for 
`)

	createAddonSentryExample = templates.Examples(`
		# Create the sentry addon 
		jx create addon sentry

		# Create the sentry addon in a custom namespace
		jx create addon sentry -n mynamespace
	`)
)

// CreateAddonSentryOptions the options for the create spring command
type CreateAddonSentryOptions struct {
	CreateAddonOptions

	Chart string
	Host  string
}

// NewCmdCreateAddonSentry creates a command object for the "create" command
func NewCmdCreateAddonSentry(f Factory, in terminal.FileReader, out terminal.FileWriter, errOut io.Writer) *cobra.Command {
	options := &CreateAddonSentryOptions{
		CreateAddonOptions: CreateAddonOptions{
			CreateOptions: CreateOptions{
				CommonOptions: CommonOptions{
					Factory: f,
					In:      in,
					Out:     out,
					Err:     errOut,
				},
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "sentry",
		Short:   "Create an sentry addon",
		Aliases: []string{"env"},
		Long:    createAddonSentryLong,
		Example: createAddonSentryExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			CheckErr(err)
		},
	}

	options.addCommonFlags(cmd)
	options.addFlags(cmd, "", defaultSentryReleaseName)

	cmd.Flags().StringVarP(&options.Version, "version", "v", "", "The version of the sentry addon to use")
	cmd.Flags().StringVarP(&options.Chart, optionChart, "c", kube.ChartSentry, "The name of the chart to use")
	cmd.Flags().StringVarP(&options.Host, "host", "", "", "The Ingress host name to expose Sentry on. Defaults to 'sentry.$namespace.$domain if ommitted")
	return cmd
}

// Run implements the command
func (o *CreateAddonSentryOptions) Run() error {
	if o.ReleaseName == "" {
		return util.MissingOption(optionRelease)
	}
	if o.Chart == "" {
		return util.MissingOption(optionChart)
	}
	err := o.ensureHelm()
	if err != nil {
		return errors.Wrap(err, "failed to ensure that helm is present")
	}
	ingressHost := o.Host
	if ingressHost == "" {
		ns := o.Namespace
		kubeClient, defaultNs, err := o.KubeClientAndDevNamespace()
		if err != nil {
			return err
		}
		if ns == "" {
			ns = defaultNs
			o.Namespace = ns
		}

		domain, err := kube.GetCurrentDomain(kubeClient, ns)
		if err != nil {
			return err
		}
		ingressHost = "sentry." + ns + "." + domain
	}

	values := strings.Split(o.SetValues, ",")
	values = helm.AddValuesIfMissing(values, "service.type=ClusterIP", "ingress.enabled=true", "ingress.hostname="+ingressHost)

	err = o.installChart(o.ReleaseName, o.Chart, o.Version, o.Namespace, true, values)
	if err != nil {
		return err
	}
	return nil
}
