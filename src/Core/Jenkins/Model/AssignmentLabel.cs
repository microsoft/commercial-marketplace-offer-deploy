using System;
using System.Text.Json.Serialization;

namespace Modm.Jenkins.Model
{
    public class AssignedLabel
    {
        [JsonPropertyName("name")]
        public string Name { get; set; }
    }
}

