using System;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using Modm.Serialization;

namespace Modm.Azure.Model
{
    public class UserData
    {
        static readonly JsonSerializerOptions serializerOptions = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            Converters = { new DictionaryStringObjectJsonConverter() }
        };

        public required string ArtifactsUri { get; set; }

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
            return !string.IsNullOrEmpty(this.ArtifactsUri)
                && Uri.IsWellFormedUriString(this.ArtifactsUri, UriKind.Absolute);
        }
    }
}

