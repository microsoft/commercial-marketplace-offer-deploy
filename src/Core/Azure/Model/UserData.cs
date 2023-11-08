using System;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using Modm.Serialization;

namespace Modm.Azure.Model
{
    public class InstallerPackageInfo
    {
        [JsonPropertyName("uri")]
        public required string Uri { get; set; }

        [JsonPropertyName("hash")]
        public required string Hash { get; set; }

        internal bool IsValid()
        {
            return !string.IsNullOrEmpty(this.Hash) &&
                   !string.IsNullOrEmpty(Uri) &&
                   System.Uri.IsWellFormedUriString(this.Uri, UriKind.Absolute);
        }
    }

    public class UserData
    {
        static readonly JsonSerializerOptions serializerOptions = new()
        {
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            Converters = { new DictionaryStringObjectJsonConverter() }
        };

        [JsonPropertyName("installerPackage")]
        public required InstallerPackageInfo InstallerPackage { get; set; }

        [JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
        public Dictionary<string, object> Parameters { get; set; }

        public string ToBase64Json()
        {
            string jsonString = JsonSerializer.Serialize(this, serializerOptions);
            byte[] jsonBytes = Encoding.UTF8.GetBytes(jsonString);

            return Convert.ToBase64String(jsonBytes);
        }

        public static UserData Deserialize(string base64UserData)
        {
            byte[] data = Convert.FromBase64String(base64UserData);
            string jsonString = Encoding.UTF8.GetString(data);

            UserData userData = JsonSerializer.Deserialize<UserData>(jsonString, serializerOptions);

            return userData;
        }

        public bool IsValid()
        {
            return this.InstallerPackage.IsValid();
        }
    }
}

