using System;
using System.Text;
using System.Text.Json;

namespace Modm.Azure.Model
{
    public class UserData
    {
        public required string ArtifactsUri { get; set; }
        public IDictionary<string, object> Properties { get; set; }

        public string ToBase64Json()
        {
            string jsonString = JsonSerializer.Serialize(this);
            byte[] jsonBytes = Encoding.UTF8.GetBytes(jsonString);

            return Convert.ToBase64String(jsonBytes);
        }

        public static UserData Deserialize(string base64UserData)
        {
            byte[] data = Convert.FromBase64String(base64UserData);
            string jsonString = Encoding.UTF8.GetString(data);

            UserData userData = JsonSerializer.Deserialize<UserData>(jsonString);

            return userData;
        }

        public bool IsValid()
        {
            return !string.IsNullOrEmpty(this.ArtifactsUri)
                && Uri.IsWellFormedUriString(this.ArtifactsUri, UriKind.Absolute);
        }
    }
}

