// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System.Collections.Generic;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace Modm.Deployments
{
    public class AuditFile : JsonFile<List<AuditRecord>>
    {
        public override string FileName => "audit.json";

        public AuditFile(IConfiguration configuration, ILogger<AuditFile> logger)
            : base(configuration, logger)
        {
        }
    }
}
