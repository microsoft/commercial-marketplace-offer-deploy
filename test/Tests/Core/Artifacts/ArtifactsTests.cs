using System;
using Mono.Unix.Native;
using Modm.Artifacts;

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
                artifactsFile.ChangeDirectoryPermissions(tempDir);

                Stat stat;

                if (Syscall.stat(tempDir, out stat) == 0)
                {
                    var mode = stat.st_mode;
                    Assert.Equal("ACCESSPERMS, S_IFDIR", mode.ToString());
                }
                else
                {
                    Assert.Fail("Could not retrieve directory stats");
                }
            }
            finally
            {
                // Clean up
                Directory.Delete(tempDir, true);
            }
        }
    }
}

