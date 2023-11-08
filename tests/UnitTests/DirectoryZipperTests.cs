using System.IO.Compression;
using Modm.Compression;

namespace Modm.Tests.UnitTests
{
    public class DirectoryZipperTests
    {
        [Fact]
        public void ZipDirectory_CreatesZipOfDirectory()
        {
            // Arrange
            var sourceDirectory = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
            var destinationZip = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString() + ".zip");
            Directory.CreateDirectory(sourceDirectory);
            File.WriteAllText(Path.Combine(sourceDirectory, "test.txt"), "Hello, World!");

            var zipper = new DirectoryZipper();

            // Act
            zipper.ZipDirectory(sourceDirectory, destinationZip);

            // Assert
            Assert.True(File.Exists(destinationZip));

            // Check if the file inside the zip is as expected
            using var archive = ZipFile.OpenRead(destinationZip);
            var entry = archive.GetEntry("test.txt");
            Assert.NotNull(entry);

            // Cleanup
            Directory.Delete(sourceDirectory, true);
            File.Delete(destinationZip);
        }

        [Fact]
        public void ZipDirectory_ThrowsWhenSourceDirectoryDoesNotExist()
        {
            // Arrange
            var nonExistentDirectory = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
            var destinationZip = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString() + ".zip");

            var zipper = new DirectoryZipper();

            // Act & Assert
            Assert.Throws<DirectoryNotFoundException>(() => zipper.ZipDirectory(nonExistentDirectory, destinationZip));
        }
    }
}