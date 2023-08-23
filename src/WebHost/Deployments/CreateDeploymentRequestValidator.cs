using System;
using FluentValidation;

namespace WebHost.Deployments
{
	public class CreateDeploymentRequestValidator : AbstractValidator<CreateDeploymentRequest>
    {
		public CreateDeploymentRequestValidator()
		{
			RuleFor(x => x.ArtifactsUri).NotEmpty().NotNull().Must(value =>
			{
				return Uri.TryCreate(value, new UriCreationOptions { DangerousDisablePathAndQueryCanonicalization = false }, out var result);
			});
		}
	}
}

