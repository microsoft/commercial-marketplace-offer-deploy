using System;
using Mono.Unix.Native;
using Modm.Artifacts;
using Mono.Unix;

namespace Modm.Tests.Core.Artifacts
{
	public class ArtifactsTests
	{
		public ArtifactsTests()
		{
		}

        [Fact]
		public void ArtifactsCanChangeDirectoryPermissions()
		{
            string tempDir = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
            Directory.CreateDirectory(tempDir);

            try
            {
                // Act
                var artifactsFile = new ArtifactsFile("");

               // Stat stat1;
                var unixFileInfo = new UnixFileInfo(tempDir);
                unixFileInfo.FileAccessPermissions = FileAccessPermissions.DefaultPermissions;

                //if (Syscall.stat(tempDir, out stat1) == 0)
                //{
                //    var mode = stat1.st_mode;
                //}

                artifactsFile.ChangeDirectoryPermissions(tempDir);
                var unixFileInfo2 = new UnixFileInfo(tempDir);
                Assert.Equal(FileAccessPermissions.AllPermissions, unixFileInfo2.FileAccessPermissions);

                //Stat stat;

                //if (Syscall.stat(tempDir, out stat) == 0)
                //{
                //    var mode = stat.st_mode;
                //    Assert.Equal("ACCESSPERMS, S_IFDIR", mode.ToString());

                    
                //}
                //else
                //{
                //    Assert.Fail("Could not retrieve directory stats");
                //}
            }
            finally
            {
                // Clean up
                Directory.Delete(tempDir, true);
            }
        }
    }
}

