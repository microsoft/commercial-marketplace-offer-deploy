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

                ArtifactsFile.ChangeDirectoryPermissions(tempDir);

                var m2 = d.UnixFileMode;
                Assert.True(m2.HasFlag(UnixFileMode.UserExecute));
            }
            finally
            {
                // Clean up
                Directory.Delete(tempDir, true);
            }
        }

        [Fact]
        public async Task signature_should_validate()
        {
            string expectedSig = "3be95f70f2a217ac0278cd62d72fdf12732f9a69e92db77f2a70f5f3e6996e2a";
            // Get the current executing directory.
            string currentDirectory = AppDomain.CurrentDomain.BaseDirectory;
            currentDirectory = Path.Combine(currentDirectory, "Core");
            currentDirectory = Path.Combine(currentDirectory, "Artifacts");

            // Construct the path to the content.zip file.
            string filePath = Path.Combine(currentDirectory, "content.zip");
            //string filePath = "./content.zip";

            var artifactsFile = new ArtifactsFile(filePath);
            bool isValid = artifactsFile.IsValidSignature(expectedSig);

            Assert.True(isValid);
        }
    }
}

