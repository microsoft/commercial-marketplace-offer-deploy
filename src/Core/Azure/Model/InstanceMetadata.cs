// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.Azure.Management.ResourceManager.Fluent.Core;

namespace Modm.Azure.Model
{
    public class InstanceMetadata
    {
        [JsonPropertyName("compute")]
        public required Compute Compute { get; set; }

        [JsonPropertyName("network")]
        public required Network Network { get; set; }

        [JsonIgnore]
        public ResourceId ResourceGroupId => ResourceId.FromString($"/subscriptions/{Compute.SubscriptionId}/resourceGroups/{Compute.ResourceGroupName}");
    }

    public partial class Compute
    {
        [JsonPropertyName("azEnvironment")]
        public string AzEnvironment { get; set; }

        [JsonPropertyName("customData")]
        public string CustomData { get; set; }

        [JsonPropertyName("evictionPolicy")]
        public string EvictionPolicy { get; set; }

        [JsonPropertyName("isHostCompatibilityLayerVm")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool IsHostCompatibilityLayerVm { get; set; }

        [JsonPropertyName("licenseType")]
        public string LicenseType { get; set; }

        [JsonPropertyName("location")]
        public string Location { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("offer")]
        public string Offer { get; set; }

        [JsonPropertyName("osProfile")]
        public OsProfile OsProfile { get; set; }

        [JsonPropertyName("osType")]
        public string OsType { get; set; }

        [JsonPropertyName("placementGroupId")]
        public string PlacementGroupId { get; set; }

        [JsonPropertyName("plan")]
        public Plan Plan { get; set; }

        [JsonPropertyName("platformFaultDomain")]
        public string PlatformFaultDomain { get; set; }

        [JsonPropertyName("platformUpdateDomain")]
        public string PlatformUpdateDomain { get; set; }

        [JsonPropertyName("priority")]
        public string Priority { get; set; }

        [JsonPropertyName("provider")]
        public string Provider { get; set; }

        [JsonPropertyName("publicKeys")]
        public PublicKey[] PublicKeys { get; set; }

        [JsonPropertyName("publisher")]
        public string Publisher { get; set; }

        [JsonPropertyName("resourceGroupName")]
        public string ResourceGroupName { get; set; }

        [JsonPropertyName("resourceId")]
        public string ResourceId { get; set; }

        [JsonPropertyName("securityProfile")]
        public SecurityProfile SecurityProfile { get; set; }

        [JsonPropertyName("sku")]
        public string Sku { get; set; }

        [JsonPropertyName("storageProfile")]
        public StorageProfile StorageProfile { get; set; }

        [JsonPropertyName("subscriptionId")]
        public Guid SubscriptionId { get; set; }

        [JsonPropertyName("tags")]
        public string Tags { get; set; }

        [JsonPropertyName("tagsList")]
        public List<KeyValuePair<string,string>> TagsList { get; set; }

        [JsonPropertyName("userData")]
        public string UserData { get; set; }

        [JsonPropertyName("version")]
        public string Version { get; set; }

        [JsonPropertyName("vmId")]
        public Guid VmId { get; set; }

        [JsonPropertyName("vmScaleSetName")]
        public string VmScaleSetName { get; set; }

        [JsonPropertyName("vmSize")]
        public string VmSize { get; set; }

        [JsonPropertyName("zone")]
        public string Zone { get; set; }
    }

    public partial class OsProfile
    {
        [JsonPropertyName("adminUsername")]
        public string AdminUsername { get; set; }

        [JsonPropertyName("computerName")]
        public string ComputerName { get; set; }

        [JsonPropertyName("disablePasswordAuthentication")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool DisablePasswordAuthentication { get; set; }
    }

    public partial class Plan
    {
        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("product")]
        public string Product { get; set; }

        [JsonPropertyName("publisher")]
        public string Publisher { get; set; }
    }

    public partial class PublicKey
    {
        [JsonPropertyName("keyData")]
        public string KeyData { get; set; }

        [JsonPropertyName("path")]
        public string Path { get; set; }
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
        public object[] DataDisks { get; set; }

        [JsonPropertyName("imageReference")]
        public ImageReference ImageReference { get; set; }

        [JsonPropertyName("osDisk")]
        public OsDisk OsDisk { get; set; }

        [JsonPropertyName("resourceDisk")]
        public ResourceDisk ResourceDisk { get; set; }
    }

    public partial class ImageReference
    {
        [JsonPropertyName("id")]
        public string Id { get; set; }

        [JsonPropertyName("offer")]
        public string Offer { get; set; }

        [JsonPropertyName("publisher")]
        public string Publisher { get; set; }

        [JsonPropertyName("sku")]
        public string Sku { get; set; }

        [JsonPropertyName("version")]
        public string Version { get; set; }
    }

    public partial class OsDisk
    {
        [JsonPropertyName("caching")]
        public string Caching { get; set; }

        [JsonPropertyName("createOption")]
        public string CreateOption { get; set; }

        [JsonPropertyName("diffDiskSettings")]
        public DiffDiskSettings DiffDiskSettings { get; set; }

        [JsonPropertyName("diskSizeGB")]
        public string DiskSizeGb { get; set; }

        [JsonPropertyName("encryptionSettings")]
        public EncryptionSettings EncryptionSettings { get; set; }

        [JsonPropertyName("image")]
        public Image Image { get; set; }

        [JsonPropertyName("managedDisk")]
        public ManagedDisk ManagedDisk { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("osType")]
        public string OsType { get; set; }

        [JsonPropertyName("vhd")]
        public Image Vhd { get; set; }

        [JsonPropertyName("writeAcceleratorEnabled")]
        [JsonConverter(typeof(BooleanConverter))]
        public bool WriteAcceleratorEnabled { get; set; }
    }

    public partial class DiffDiskSettings
    {
        [JsonPropertyName("option")]
        public string Option { get; set; }
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
        public string Uri { get; set; }
    }

    public partial class ManagedDisk
    {
        [JsonPropertyName("id")]
        public string Id { get; set; }

        [JsonPropertyName("storageAccountType")]
        public string StorageAccountType { get; set; }
    }

    public partial class ResourceDisk
    {
        [JsonPropertyName("size")]
        public string Size { get; set; }
    }

    public partial class Network
    {
        [JsonPropertyName("interface")]
        public Interface[] Interface { get; set; }
    }

    public partial class Interface
    {
        [JsonPropertyName("ipv4")]
        public Ipv4 Ipv4 { get; set; }

        [JsonPropertyName("ipv6")]
        public Ipv6 Ipv6 { get; set; }

        [JsonPropertyName("macAddress")]
        public string MacAddress { get; set; }
    }

    public partial class Ipv4
    {
        [JsonPropertyName("ipAddress")]
        public IpAddress[] IpAddress { get; set; }

        [JsonPropertyName("subnet")]
        public Subnet[] Subnet { get; set; }
    }

    public partial class IpAddress
    {
        [JsonPropertyName("privateIpAddress")]
        public string PrivateIpAddress { get; set; }

        [JsonPropertyName("publicIpAddress")]
        public string PublicIpAddress { get; set; }
    }

    public partial class Subnet
    {
        [JsonPropertyName("address")]
        public string Address { get; set; }

        [JsonPropertyName("prefix")]
        public string Prefix { get; set; }
    }

    public partial class Ipv6
    {
        [JsonPropertyName("ipAddress")]
        public object[] IpAddress { get; set; }
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
