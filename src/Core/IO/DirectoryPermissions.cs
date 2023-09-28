using System;
namespace Modm.IO
{
	public class DirectoryPermissions
	{
        /// <summary>
        /// This method will change the permissions of the directory to allow everyone to have full control
        /// </summary>
        /// <param name="directoryPath"></param>
        /// <exception cref="DirectoryNotFoundException"></exception>
        [System.Diagnostics.CodeAnalysis.SuppressMessage("Interoperability",
            "CA1416:Validate platform compatibility",
            Justification = "We are only targeting *Nix OS, never windows")]
        public static void AllowFullControl(string directoryPath)
        {
            try
            {
                DirectoryInfo directoryInfo = new(directoryPath);

                var ownerPermissions = UnixFileMode.UserExecute | UnixFileMode.UserRead | UnixFileMode.UserWrite;
                var groupPermissions = UnixFileMode.GroupExecute | UnixFileMode.GroupRead | UnixFileMode.GroupWrite;
                var othersPermissions = UnixFileMode.OtherExecute | UnixFileMode.OtherRead | UnixFileMode.OtherWrite;
                var allPermissions = ownerPermissions | groupPermissions | othersPermissions;

                directoryInfo.UnixFileMode = allPermissions;
            }
            catch (Exception)
            {

            }
        }
    }
}

