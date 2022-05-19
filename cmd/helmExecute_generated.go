// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type helmExecuteOptions struct {
	AdditionalParameters      []string `json:"additionalParameters,omitempty"`
	ChartPath                 string   `json:"chartPath,omitempty"`
	TargetRepositoryURL       string   `json:"targetRepositoryURL,omitempty"`
	TargetRepositoryName      string   `json:"targetRepositoryName,omitempty"`
	TargetRepositoryUser      string   `json:"targetRepositoryUser,omitempty"`
	TargetRepositoryPassword  string   `json:"targetRepositoryPassword,omitempty"`
	HelmDeployWaitSeconds     int      `json:"helmDeployWaitSeconds,omitempty"`
	HelmValues                []string `json:"helmValues,omitempty"`
	Image                     string   `json:"image,omitempty"`
	KeepFailedDeployments     bool     `json:"keepFailedDeployments,omitempty"`
	KubeConfig                string   `json:"kubeConfig,omitempty"`
	KubeContext               string   `json:"kubeContext,omitempty"`
	Namespace                 string   `json:"namespace,omitempty"`
	DockerConfigJSON          string   `json:"dockerConfigJSON,omitempty"`
	HelmCommand               string   `json:"helmCommand,omitempty" validate:"possible-values=upgrade lint install test uninstall dependency publish"`
	AppVersion                string   `json:"appVersion,omitempty"`
	Dependency                string   `json:"dependency,omitempty" validate:"possible-values=build list update"`
	PackageDependencyUpdate   bool     `json:"packageDependencyUpdate,omitempty"`
	DumpLogs                  bool     `json:"dumpLogs,omitempty"`
	FilterTest                string   `json:"filterTest,omitempty"`
	CustomTLSCertificateLinks []string `json:"customTlsCertificateLinks,omitempty"`
	Publish                   bool     `json:"publish,omitempty"`
	Version                   string   `json:"version,omitempty"`
}

// HelmExecuteCommand Executes helm3 functionality as the package manager for Kubernetes.
func HelmExecuteCommand() *cobra.Command {
	const STEP_NAME = "helmExecute"

	metadata := helmExecuteMetadata()
	var stepConfig helmExecuteOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createHelmExecuteCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Executes helm3 functionality as the package manager for Kubernetes.",
		Long: `Alpha version: please expect incompatible changes

Executes helm functionality as the package manager for Kubernetes.

* [Helm](https://helm.sh/)  is the package manager for Kubernetes.
* [Helm documentation https://helm.sh/docs/intro/using_helm/ and best practies https://helm.sh/docs/chart_best_practices/conventions/]
* [Helm Charts] (https://artifacthub.io/)
` + "`" + `` + "`" + `` + "`" + `
Available Commands:
` + "`" + `upgrade` + "`" + `, ` + "`" + `lint` + "`" + `, ` + "`" + `install` + "`" + `, ` + "`" + `test` + "`" + `, ` + "`" + `uninstall` + "`" + `, ` + "`" + `dependency` + "`" + `, ` + "`" + `publish` + "`" + `

  upgrade       upgrade a release
  lint          examine a chart for possible issues
  install       install a chart
  test          run tests for a release
  uninstall     uninstall a release
  dependency     package a chart directory into a chart archive
  publish       package and puslish a release

` + "`" + `` + "`" + `` + "`" + `

Note: piper supports only helm3 version, since helm2 is deprecated.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.TargetRepositoryUser)
			log.RegisterSecret(stepConfig.TargetRepositoryPassword)
			log.RegisterSecret(stepConfig.KubeConfig)
			log.RegisterSecret(stepConfig.DockerConfigJSON)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			helmExecute(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addHelmExecuteFlags(createHelmExecuteCmd, &stepConfig)
	return createHelmExecuteCmd
}

func addHelmExecuteFlags(cmd *cobra.Command, stepConfig *helmExecuteOptions) {
	cmd.Flags().StringSliceVar(&stepConfig.AdditionalParameters, "additionalParameters", []string{}, "Defines additional parameters for Helm like  \"helm install [NAME] [CHART] [flags]\".")
	cmd.Flags().StringVar(&stepConfig.ChartPath, "chartPath", os.Getenv("PIPER_chartPath"), "Defines the chart path for helm. chartPath is mandatory for install/upgrade/publish commands.")
	cmd.Flags().StringVar(&stepConfig.TargetRepositoryURL, "targetRepositoryURL", os.Getenv("PIPER_targetRepositoryURL"), "URL of the target repository where the compiled helm .tgz archive shall be uploaded - typically provided by the CI/CD environment.")
	cmd.Flags().StringVar(&stepConfig.TargetRepositoryName, "targetRepositoryName", os.Getenv("PIPER_targetRepositoryName"), "set the chart repository. The value is required for install/upgrade/uninstall commands.")
	cmd.Flags().StringVar(&stepConfig.TargetRepositoryUser, "targetRepositoryUser", os.Getenv("PIPER_targetRepositoryUser"), "Username for the char repository where the compiled helm .tgz archive shall be uploaded - typically provided by the CI/CD environment.")
	cmd.Flags().StringVar(&stepConfig.TargetRepositoryPassword, "targetRepositoryPassword", os.Getenv("PIPER_targetRepositoryPassword"), "Password for the target repository where the compiled helm .tgz archive shall be uploaded - typically provided by the CI/CD environment.")
	cmd.Flags().IntVar(&stepConfig.HelmDeployWaitSeconds, "helmDeployWaitSeconds", 300, "Number of seconds before helm deploy returns.")
	cmd.Flags().StringSliceVar(&stepConfig.HelmValues, "helmValues", []string{}, "List of helm values as YAML file reference or URL (as per helm parameter description for `-f` / `--values`)")
	cmd.Flags().StringVar(&stepConfig.Image, "image", os.Getenv("PIPER_image"), "Full name of the image to be deployed.")
	cmd.Flags().BoolVar(&stepConfig.KeepFailedDeployments, "keepFailedDeployments", false, "Defines whether a failed deployment will be purged")
	cmd.Flags().StringVar(&stepConfig.KubeConfig, "kubeConfig", os.Getenv("PIPER_kubeConfig"), "Defines the path to the \"kubeconfig\" file.")
	cmd.Flags().StringVar(&stepConfig.KubeContext, "kubeContext", os.Getenv("PIPER_kubeContext"), "Defines the context to use from the \"kubeconfig\" file.")
	cmd.Flags().StringVar(&stepConfig.Namespace, "namespace", `default`, "Defines the target Kubernetes namespace for the deployment.")
	cmd.Flags().StringVar(&stepConfig.DockerConfigJSON, "dockerConfigJSON", os.Getenv("PIPER_dockerConfigJSON"), "Path to the file `.docker/config.json` - this is typically provided by your CI/CD system. You can find more details about the Docker credentials in the [Docker documentation](https://docs.docker.com/engine/reference/commandline/login/).")
	cmd.Flags().StringVar(&stepConfig.HelmCommand, "helmCommand", os.Getenv("PIPER_helmCommand"), "Helm: defines the command `upgrade`, `lint`, `install`, `test`, `uninstall`, `dependency`, `publish`.")
	cmd.Flags().StringVar(&stepConfig.AppVersion, "appVersion", os.Getenv("PIPER_appVersion"), "set the appVersion on the chart to this version")
	cmd.Flags().StringVar(&stepConfig.Dependency, "dependency", os.Getenv("PIPER_dependency"), "manage a chart's dependencies")
	cmd.Flags().BoolVar(&stepConfig.PackageDependencyUpdate, "packageDependencyUpdate", false, "update dependencies from \"Chart.yaml\" to dir \"charts/\" before packaging")
	cmd.Flags().BoolVar(&stepConfig.DumpLogs, "dumpLogs", false, "dump the logs from test pods (this runs after all tests are complete, but before any cleanup)")
	cmd.Flags().StringVar(&stepConfig.FilterTest, "filterTest", os.Getenv("PIPER_filterTest"), "specify tests by attribute (currently `name`) using attribute=value syntax or `!attribute=value` to exclude a test (can specify multiple or separate values with commas `name=test1,name=test2`)")
	cmd.Flags().StringSliceVar(&stepConfig.CustomTLSCertificateLinks, "customTlsCertificateLinks", []string{}, "List of download links to custom TLS certificates. This is required to ensure trusted connections to instances with repositories (like nexus) when publish flag is set to true.")
	cmd.Flags().BoolVar(&stepConfig.Publish, "publish", false, "Configures helm to run the deploy command to publish artifacts to a repository.")
	cmd.Flags().StringVar(&stepConfig.Version, "version", os.Getenv("PIPER_version"), "Defines the artifact version to use from helm package/publish commands.")

	cmd.MarkFlagRequired("image")
}

// retrieve step metadata
func helmExecuteMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "helmExecute",
			Aliases:     []config.Alias{},
			Description: "Executes helm3 functionality as the package manager for Kubernetes.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "dockerCredentialsId", Type: "jenkins"},
					{Name: "dockerConfigJsonCredentialsId", Description: "Jenkins 'Secret file' credentials ID containing Docker config.json (with registry credential(s)).", Type: "jenkins"},
				},
				Resources: []config.StepResources{
					{Name: "deployDescriptor", Type: "stash"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "additionalParameters",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmDeploymentParameters"}},
						Default:     []string{},
					},
					{
						Name:        "chartPath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmChartPath"}},
						Default:     os.Getenv("PIPER_chartPath"),
					},
					{
						Name: "targetRepositoryURL",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/helmRepositoryURL",
							},

							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryUrl",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_targetRepositoryURL"),
					},
					{
						Name:        "targetRepositoryName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_targetRepositoryName"),
					},
					{
						Name: "targetRepositoryUser",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/helmRepositoryUsername",
							},

							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryUsername",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_targetRepositoryUser"),
					},
					{
						Name: "targetRepositoryPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/helmRepositoryPassword",
							},

							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryPassword",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_targetRepositoryPassword"),
					},
					{
						Name:        "helmDeployWaitSeconds",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "int",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     300,
					},
					{
						Name:        "helmValues",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name: "image",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "container/imageNameTag",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "deployImage"}},
						Default:   os.Getenv("PIPER_image"),
					},
					{
						Name:        "keepFailedDeployments",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name: "kubeConfig",
						ResourceRef: []config.ResourceReference{
							{
								Name: "kubeConfigFileCredentialsId",
								Type: "secret",
							},

							{
								Name:    "kubeConfigFileVaultSecretName",
								Type:    "vaultSecretFile",
								Default: "kube-config",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_kubeConfig"),
					},
					{
						Name:        "kubeContext",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_kubeContext"),
					},
					{
						Name:        "namespace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmDeploymentNamespace"}},
						Default:     `default`,
					},
					{
						Name: "dockerConfigJSON",
						ResourceRef: []config.ResourceReference{
							{
								Name: "dockerConfigJsonCredentialsId",
								Type: "secret",
							},

							{
								Name:    "dockerConfigFileVaultSecretName",
								Type:    "vaultSecretFile",
								Default: "docker-config",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_dockerConfigJSON"),
					},
					{
						Name:        "helmCommand",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_helmCommand"),
					},
					{
						Name:        "appVersion",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_appVersion"),
					},
					{
						Name:        "dependency",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_dependency"),
					},
					{
						Name:        "packageDependencyUpdate",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "dumpLogs",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "filterTest",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_filterTest"),
					},
					{
						Name:        "customTlsCertificateLinks",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "publish",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "version",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_version"),
					},
				},
			},
			Containers: []config.Container{
				{Image: "dtzar/helm-kubectl:3.8.0", WorkingDir: "/config", Options: []config.Option{{Name: "-u", Value: "0"}}},
			},
		},
	}
	return theMetaData
}
