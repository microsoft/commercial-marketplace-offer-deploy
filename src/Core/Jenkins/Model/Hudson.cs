// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json.Serialization;

namespace Modm.Jenkins.Model
{
    /// <summary>
    /// Generated from JSON from http://localhost:8080/api/json
    /// Provides information about Jenkins
    /// </summary>
    public class Hudson
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("assignedLabels")]
        public List<AssignedLabel> AssignedLabels { get; set; }

        [JsonPropertyName("mode")]
        public string Mode { get; set; }

        [JsonPropertyName("nodeDescription")]
        public string NodeDescription { get; set; }

        [JsonPropertyName("nodeName")]
        public string NodeName { get; set; }

        [JsonPropertyName("numExecutors")]
        public int NumExecutors { get; set; }

        [JsonPropertyName("description")]
        public object Description { get; set; }

        [JsonPropertyName("jobs")]
        public List<Job> Jobs { get; set; }

        [JsonPropertyName("overallLoad")]
        public OverallLoad OverallLoad { get; set; }

        [JsonPropertyName("primaryView")]
        public PrimaryView PrimaryView { get; set; }

        [JsonPropertyName("quietDownReason")]
        public object QuietDownReason { get; set; }

        [JsonPropertyName("quietingDown")]
        public bool QuietingDown { get; set; }

        [JsonPropertyName("slaveAgentPort")]
        public int SlaveAgentPort { get; set; }

        [JsonPropertyName("unlabeledLoad")]
        public UnlabeledLoad UnlabeledLoad { get; set; }

        [JsonPropertyName("url")]
        public object Url { get; set; }

        [JsonPropertyName("useCrumbs")]
        public bool UseCrumbs { get; set; }

        [JsonPropertyName("useSecurity")]
        public bool UseSecurity { get; set; }

        [JsonPropertyName("views")]
        public List<View> Views { get; set; }
    }

    public class Job
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("url")]
        public string Url { get; set; }

        [JsonPropertyName("color")]
        public string Color { get; set; }
    }

    public class OverallLoad
    {
    }

    public class PrimaryView
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("url")]
        public string Url { get; set; }
    }

    public class UnlabeledLoad
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }
    }

    public class View
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("url")]
        public string Url { get; set; }
    }


}

