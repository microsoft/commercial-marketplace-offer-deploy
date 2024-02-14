// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm.Deployments
{
    /// <summary>
    /// Defines the types of deployments that MODM supports
    /// </summary>
    public readonly struct DeploymentType
	{
		public static readonly string Arm = "arm";
        public static readonly string Terraform = "terraform";

        public static readonly List<string> SupportedTypes = new() { Terraform, Arm };

        private readonly string deploymentType;

        public DeploymentType(string deploymentType)
        {
            if (string.IsNullOrEmpty(deploymentType))
            {
                throw new ArgumentNullException(nameof(deploymentType));
            }

            if (deploymentType != Arm && deploymentType != Terraform)
            {
                throw new ArgumentOutOfRangeException(nameof(deploymentType), "deployment type does not match any support type.");
            }
            this.deploymentType = deploymentType;
        }

        public static implicit operator string(DeploymentType t) => t.deploymentType;
        public static implicit operator DeploymentType(string t) => new(t);

        public override readonly string ToString() => $"{deploymentType}";
    }
}

