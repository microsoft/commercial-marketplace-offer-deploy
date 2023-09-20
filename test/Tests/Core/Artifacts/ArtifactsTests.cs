using System;
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
            var d = Directory.CreateDirectory(tempDir);
           
            try
            {
                // Act
                var artifactsFile = new ArtifactsFile("");

                var m = d.UnixFileMode;
                d.UnixFileMode = UnixFileMode.OtherWrite | UnixFileMode.OtherRead | UnixFileMode.GroupRead | UnixFileMode.UserWrite | UnixFileMode.UserRead;

                artifactsFile.ChangeDirectoryPermissions(tempDir);

                var m2 = d.UnixFileMode;
                Assert.True(m2.HasFlag(UnixFileMode.UserExecute));
            }
            finally
            {
                // Clean up
                Directory.Delete(tempDir, true);
            }
        }
    }
}

