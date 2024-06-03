locals {
  private_dns_zone_id = (
    var.is_gov ?
    "/subscriptions//resourceGroups//providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.usgovcloudapi.net" :
    "/subscriptions//resourceGroups//providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.azure.com"
  )
}
