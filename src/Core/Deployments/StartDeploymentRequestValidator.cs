using FluentValidation;

namespace Modm.Deployments
{
    public class StartDeploymentRequestValidator : AbstractValidator<StartDeploymentRequest>
    {
		public StartDeploymentRequestValidator()
		{
			RuleFor(x => x.ArtifactsUri).NotEmpty().NotNull().Must(value =>
			{
				return Uri.TryCreate(value, new UriCreationOptions { DangerousDisablePathAndQueryCanonicalization = false }, out var result);
			});

			RuleFor(x => x.ArtifactsSig).NotEmpty().NotNull();

			RuleFor(x => x.Parameters).NotNull();
		} 
	}
}

