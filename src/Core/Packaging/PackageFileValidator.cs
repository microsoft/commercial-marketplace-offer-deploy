// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using FluentValidation;

namespace Modm.Packaging
{
	public class PackageFileValidator : AbstractValidator<PackageFile>
	{
		public PackageFileValidator()
		{
			RuleFor(f => f.ExtractedTo).NotEmpty().NotNull();
			RuleFor(f => f.ComputedHash).Custom((hash, context) =>
			{
				var compareTo = context.RootContextData[PackageFile.HashAttributeName];

                if (!hash.Equals(compareTo))
				{
					context.AddFailure("Installer package hash values do not match.");
				}
			});
		}
	}
}

