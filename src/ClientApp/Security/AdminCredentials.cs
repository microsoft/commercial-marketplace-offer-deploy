namespace ClientApp.Security
{
    /// <summary>
    /// Admin credentials that are set from Installer step in the createUiDefinition.json
    /// </summary>
    public class AdminCredentials
    {
        public string Username { get; set; }
        public string Password { get; set; }
    }
}