using Modm.IO;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class DirectoryPermissionsTests : IDisposable
	{
        readonly DisposableDirectory<DirectoryPermissionsTests> directory;

        public DirectoryPermissionsTests()
        {
            directory = Test.Directory<DirectoryPermissionsTests>();
        }

        [Fact]
        public void Changes_Directory_Permissions_To_Full_Control()
        {
            DirectoryPermissions.AllowFullControl(directory.FullName);
            Assert.True(directory.Info.UnixFileMode.HasFlag(UnixFileMode.UserExecute));
        }

        // cleanup
        public void Dispose()
        {
            directory.Dispose();
        }
    }
}

