using System;
using System.Text;
using Newtonsoft.Json;
namespace Modm.Azure.Model
{
    public class UserData
    {
        public required string ArtifactsUri { get; set; }

        public string ToBase64Json()
        {
            string jsonString = JsonConvert.SerializeObject(this);
            byte[] jsonBytes = Encoding.UTF8.GetBytes(jsonString);

            return Convert.ToBase64String(jsonBytes);
        }

        public bool IsValid()
        {
            return !string.IsNullOrEmpty(this.ArtifactsUri)
                && Uri.IsWellFormedUriString(this.ArtifactsUri, UriKind.Absolute);
        }
    }
}

