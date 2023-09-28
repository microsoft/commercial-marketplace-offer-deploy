using System;
using FluentValidation;

namespace Modm.Artifacts
{
	public class ArtifactsFileValidator : AbstractValidator<ArtifactsFile>
	{
		public ArtifactsFileValidator()
		{
			RuleFor(f => f.ExtractedTo).NotEmpty().NotNull();
			RuleFor(f => f.ComputedHash).Custom((hash, context) =>
			{
				var compareTo = context.RootContextData[ArtifactsFile.HashAttributeName];

                if (!hash.Equals(compareTo))
				{
					context.AddFailure("Artifacts hash values do not match.");
				}
			});
		}
	}
}

