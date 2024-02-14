// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.Text.Json;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Modm.Extensions;

namespace Modm.Deployments
{
    using Microsoft.Extensions.Configuration;
    using Microsoft.Extensions.Logging;

    
    public class DeploymentFile : JsonFile<Deployment>
    {
        public override string FileName => "deployment.json";

        public DeploymentFile(IConfiguration configuration, ILogger<DeploymentFile> logger)
            : base(configuration, logger)
        {
        }
    }
}

