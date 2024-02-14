// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Modm.Deployments
{
    public class AuditRecord
    {
        [JsonExtensionData]
        public Dictionary<string, object> AdditionalData { get; set; } = new Dictionary<string, object>();

        public AuditRecord()
        {
            AdditionalData.Add("timestamp", DateTimeOffset.UtcNow);
        }
    }

    public class DeploymentRecord
    {
        public Deployment Deployment { get; set; }
        public List<AuditRecord> AuditRecords { get; set; } = new List<AuditRecord>();

        public DeploymentRecord(Deployment deployment)
        {
            Deployment = deployment;
        }
    }
}

