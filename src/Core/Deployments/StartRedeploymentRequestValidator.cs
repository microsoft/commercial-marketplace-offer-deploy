using FluentValidation;

namespace Modm.Deployments
{
    public class StartRedeploymentRequestValidator : AbstractValidator<StartRedeploymentRequest>
    {
        public StartRedeploymentRequestValidator()
        {
            RuleFor(x => x.Parameters).NotNull();
            RuleFor(x => x.DeploymentId).NotEmpty();
        }
    }
}

