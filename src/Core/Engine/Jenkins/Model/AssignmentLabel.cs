using System;
using System.Text.Json.Serialization;

namespace Modm.Engine.Jenkins.Model
{
    class AssignedLabel
    {
        [JsonPropertyName("name")]
        public string Name { get; set; }
    }
}

