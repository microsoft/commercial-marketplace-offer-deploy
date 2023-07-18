package deployment

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type TerraformDeployer struct {
	templateType DeploymentType
}

func (deployer *TerraformDeployer) Type() DeploymentType {
	deployer.templateType = DeploymentTypeTerraform
	return deployer.templateType
}

func (deployer *TerraformDeployer) Begin(ctx context.Context, terraformDeployment TerraformDeployment) (*BeginTerraformDeploymentResult, error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}
	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	// workingDir := "/path/to/working/dir"
	workingDir := terraformDeployment.WorkingDirectory

	// save the template to the working directory
	// save the parameters file to the working directory
	// save the state file to the working directory

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}
	log.Tracef("Terraform version: %s", state.FormatVersion)
	return nil, nil
}

