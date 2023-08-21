$clientId = "<Client ID>"
$client_secret = "<Client Secret>"
$tenantId = "<Tenant ID>"

$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "application/x-www-form-urlencoded")
$body = "grant_type=client_credentials&client_id=$clientId&client_secret=$clientSecret&resource=https%3A%2F%2Fmanagement.core.windows.net"
$response = Invoke-RestMethod "https://login.microsoftonline.com/$tenantId/oauth2/token" -Method 'POST' -Headers $headers -Body $body
$accessToken = $response.access_token

$accesstoken = "token"
$headers = @{ "Authorization" = "Bearer $accesstoken" }
$DNSName = "<Machine DNS Name>"
$UserName = "<User ID>"
$Password = "<Password>"
$OS = "Linux"
$PortNo = "22"
$CompanyName = "ABCD"
$AppId = "<Application ID>"
$TenantId = "<Tenant ID>"

$body = @{
   "DNSName" = $DNSName
   "UserName" = $UserName
   "Password" = $Password
   "OS" = $OS
   "PortNo" = $PortNo
   "CompanyName" = $CompanyName
   "AppId" = $AppId
   "TenantId" = $TenantId
} | ConvertTo-Json

$body

$uri = "https://isvapp.azurewebsites.net/selftest-vm"

$res = (Invoke-WebRequest -Method "Post" -Uri $uri -Body $body -ContentType "application/json" -Headers $headers).Content