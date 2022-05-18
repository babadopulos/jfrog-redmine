package commands

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/generic"
	xrutils "github.com/jfrog/jfrog-cli-core/v2/xray/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"jfrog-redmine/service"
	"os"
)

func GetRedmineCommand() components.Command {
	return components.Command{
		Name:        "audit",
		Description: "Audit mvn project",
		Aliases:     []string{"a"},
		Arguments:   getRedmineArguments(),
		Flags:       getRedmineFlags(),
		EnvVars:     getRedmineEnvVar(),
		Action: func(c *components.Context) error {
			return redmineCmd(c)
		},
	}
}

func getRedmineArguments() []components.Argument {
	return []components.Argument{}
}

func getRedmineFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:        "source",
			Description: "Source code directory.",
			Mandatory:   true,
		},
		components.StringFlag{
			Name:        "project",
			Description: "Redmine Project identifier",
			Mandatory:   true,
		},
		components.BoolFlag{
			Name:         "dryrun",
			Description:  "Show what would have been done",
			DefaultValue: false,
		},
	}
}

func getRedmineEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        "REDMINE_API_ENDPOINT",
			Default:     "",
			Description: "Redmine API endpoint",
		},
		{
			Name:        "REDMINE_API_KEY",
			Default:     "",
			Description: "Redmine API key",
		},
	}
}

func redmineCmd(c *components.Context) error {

	var conf = new(service.RedmineConfiguration)
	conf.Source = c.GetStringFlagValue("source")
	conf.Project = c.GetStringFlagValue("project")
	conf.DryRun = c.GetBoolFlagValue("dryrun")
	conf.APIEndpoint = os.Getenv("REDMINE_API_ENDPOINT")
	conf.APIKey = os.Getenv("REDMINE_API_KEY")

	if conf.APIKey == "" {
		log.Error("Env REDMINE_API_KEY must be set")
	}

	if conf.APIKey == "" {
		log.Error("Env REDMINE_API_KEY must be set")
	}

	results, err := AuditSourceCode(c, conf)

	if err != nil {
		log.Error(err)
		return nil
	}

	issues := service.GetIssues(conf)
	service.MergeIssues(conf, results, issues)

	return nil
}

func AuditSourceCode(c *components.Context, conf *service.RedmineConfiguration) (results []services.ScanResponse, err error) {

	// Audit command works on current directory
	os.Chdir(conf.Source)

	serverConf, err := config.GetSpecificConfig("", true, false)

	auditCommand := audit.NewGenericAuditCommand()
	auditCommand.SetServerDetails(serverConf)

	technologies := []string{coreutils.Maven}

	auditCommand.SetIncludeVulnerabilities(true)
	auditCommand.SetTechnologies(technologies)
	auditCommand.SetOutputFormat(xrutils.SimpleJson)

	args := []string{}

	results, _, err = audit.GenericAudit(auditCommand.CreateXrayGraphScanParams(), serverConf, false, false, false, args, technologies...)

	return
}
