using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Modm.Azure
{
    public class InstanceMetadata
    {
        [JsonPropertyName("compute")]
        public required Compute Compute { get; set; }

        [JsonPropertyName("network")]
        public required Network Network { get; set; }
    }

    public partial class Compute
    {
        [JsonPropertyName("azEnvironment")]
        public required string AzEnvironment { get; set; }

        [JsonPropertyName("customData")]
        public required string CustomData { get; set; }

        [JsonPropertyName("evictionPolicy")]
        public required string EvictionPolicy { get; set; }

        [JsonPropertyName("isHostCompatibilityLayerVm")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool IsHostCompatibilityLayerVm { get; set; }

        [JsonPropertyName("licenseType")]
        public required string LicenseType { get; set; }

        [JsonPropertyName("location")]
        public required string Location { get; set; }

        [JsonPropertyName("name")]
        public required string Name { get; set; }

        [JsonPropertyName("offer")]
        public required string Offer { get; set; }

        [JsonPropertyName("osProfile")]
        public required OsProfile OsProfile { get; set; }

        [JsonPropertyName("osType")]
        public required string OsType { get; set; }

        [JsonPropertyName("placementGroupId")]
        public required string PlacementGroupId { get; set; }

        [JsonPropertyName("plan")]
        public required Plan Plan { get; set; }

        [JsonPropertyName("platformFaultDomain")]
        public long PlatformFaultDomain { get; set; }

        [JsonPropertyName("platformUpdateDomain")]
        public long PlatformUpdateDomain { get; set; }

        [JsonPropertyName("priority")]
        public required string Priority { get; set; }

        [JsonPropertyName("provider")]
        public required string Provider { get; set; }

        [JsonPropertyName("publicKeys")]
        public PublicKey[]? PublicKeys { get; set; }

        [JsonPropertyName("publisher")]
        public required string Publisher { get; set; }

        [JsonPropertyName("resourceGroupName")]
        public required string ResourceGroupName { get; set; }

        [JsonPropertyName("resourceId")]
        public required string ResourceId { get; set; }

        [JsonPropertyName("securityProfile")]
        public SecurityProfile? SecurityProfile { get; set; }

        [JsonPropertyName("sku")]
        public required string Sku { get; set; }

        [JsonPropertyName("storageProfile")]
        public required StorageProfile StorageProfile { get; set; }

        [JsonPropertyName("subscriptionId")]
        public Guid SubscriptionId { get; set; }

        [JsonPropertyName("tags")]
        public required string Tags { get; set; }

        [JsonPropertyName("tagsList")]
        public required List<KeyValuePair<string,string>> TagsList { get; set; }

        [JsonPropertyName("userData")]
        public required string UserData { get; set; }

        [JsonPropertyName("version")]
        public required string Version { get; set; }

        [JsonPropertyName("vmId")]
        public Guid VmId { get; set; }

        [JsonPropertyName("vmScaleSetName")]
        public required string VmScaleSetName { get; set; }

        [JsonPropertyName("vmSize")]
        public required string VmSize { get; set; }

        [JsonPropertyName("zone")]
        public required string Zone { get; set; }
    }

    public partial class OsProfile
    {
        [JsonPropertyName("adminUsername")]
        public required string AdminUsername { get; set; }

        [JsonPropertyName("computerName")]
        public required string ComputerName { get; set; }

        [JsonPropertyName("disablePasswordAuthentication")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool DisablePasswordAuthentication { get; set; }
    }

    public partial class Plan
    {
        [JsonPropertyName("name")]
        public required string Name { get; set; }

        [JsonPropertyName("product")]
        public required string Product { get; set; }

        [JsonPropertyName("publisher")]
        public required string Publisher { get; set; }
    }

    public partial class PublicKey
    {
        [JsonPropertyName("keyData")]
        public required string KeyData { get; set; }

        [JsonPropertyName("path")]
        public required string Path { get; set; }
    }

    public partial class SecurityProfile
    {
        [JsonPropertyName("secureBootEnabled")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool SecureBootEnabled { get; set; }

        [JsonPropertyName("virtualTpmEnabled")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool VirtualTpmEnabled { get; set; }
    }

    public class StorageProfile
    {
        [JsonPropertyName("dataDisks")]
        public required object[] DataDisks { get; set; }

        [JsonPropertyName("imageReference")]
        public required ImageReference ImageReference { get; set; }

        [JsonPropertyName("osDisk")]
        public required OsDisk OsDisk { get; set; }

        [JsonPropertyName("resourceDisk")]
        public ResourceDisk? ResourceDisk { get; set; }
    }

    public partial class ImageReference
    {
        [JsonPropertyName("id")]
        public required string Id { get; set; }

        [JsonPropertyName("offer")]
        public required string Offer { get; set; }

        [JsonPropertyName("publisher")]
        public required string Publisher { get; set; }

        [JsonPropertyName("sku")]
        public required string Sku { get; set; }

        [JsonPropertyName("version")]
        public required string Version { get; set; }
    }

    public partial class OsDisk
    {
        [JsonPropertyName("caching")]
        public required string Caching { get; set; }

        [JsonPropertyName("createOption")]
        public required string CreateOption { get; set; }

        [JsonPropertyName("diffDiskSettings")]
        public DiffDiskSettings? DiffDiskSettings { get; set; }

        [JsonPropertyName("diskSizeGB")]
        public long DiskSizeGb { get; set; }

        [JsonPropertyName("encryptionSettings")]
        public EncryptionSettings? EncryptionSettings { get; set; }

        [JsonPropertyName("image")]
        public Image? Image { get; set; }

        [JsonPropertyName("managedDisk")]
        public ManagedDisk? ManagedDisk { get; set; }

        [JsonPropertyName("name")]
        public required string Name { get; set; }

        [JsonPropertyName("osType")]
        public required string OsType { get; set; }

        [JsonPropertyName("vhd")]
        public Image? Vhd { get; set; }

        [JsonPropertyName("writeAcceleratorEnabled")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool WriteAcceleratorEnabled { get; set; }
    }

    public partial class DiffDiskSettings
    {
        [JsonPropertyName("option")]
        public required string Option { get; set; }
    }

    public partial class EncryptionSettings
    {
        [JsonPropertyName("enabled")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool Enabled { get; set; }
    }

    public partial class Image
    {
        [JsonPropertyName("uri")]
        public required string Uri { get; set; }
    }

    public partial class ManagedDisk
    {
        [JsonPropertyName("id")]
        public required string Id { get; set; }

        [JsonPropertyName("storageAccountType")]
        public required string StorageAccountType { get; set; }
    }

    public partial class ResourceDisk
    {
        [JsonPropertyName("size")]
        public long Size { get; set; }
    }

    public partial class Network
    {
        [JsonPropertyName("interface")]
        public Interface[]? Interface { get; set; }
    }

    public partial class Interface
    {
        [JsonPropertyName("ipv4")]
        public required Ipv4 Ipv4 { get; set; }

        [JsonPropertyName("ipv6")]
        public required Ipv6 Ipv6 { get; set; }

        [JsonPropertyName("macAddress")]
        public required string MacAddress { get; set; }
    }

    public partial class Ipv4
    {
        [JsonPropertyName("ipAddress")]
        public IpAddress[]? IpAddress { get; set; }

        [JsonPropertyName("subnet")]
        public required Subnet[] Subnet { get; set; }
    }

    public partial class IpAddress
    {
        [JsonPropertyName("privateIpAddress")]
        public required string PrivateIpAddress { get; set; }

        [JsonPropertyName("publicIpAddress")]
        public required string PublicIpAddress { get; set; }
    }

    public partial class Subnet
    {
        [JsonPropertyName("address")]
        public required string Address { get; set; }

        [JsonPropertyName("prefix")]
        public long Prefix { get; set; }
    }

    public partial class Ipv6
    {
        [JsonPropertyName("ipAddress")]
        public required object[] IpAddress { get; set; }
    }

    public class BooleanConverter : JsonConverter<bool>
    {
        public override bool Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            switch (reader.TokenType)
            {
                case JsonTokenType.True:
                    return true;
                case JsonTokenType.False:
                    return false;
                case JsonTokenType.String:
                    return reader.GetString() switch
                    {
                        "true" => true,
                        "false" => false,
                        _ => throw new JsonException()
                    };
                default:
                    throw new JsonException();
            }
        }

        public override void Write(Utf8JsonWriter writer, bool value, JsonSerializerOptions options)
        {
            writer.WriteBooleanValue(value);
        }
    }
}
