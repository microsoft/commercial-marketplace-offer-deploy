using Microsoft.Extensions.DependencyInjection;
using Modm.Packaging;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class PackageFileTests : AbstractTest<PackageFileTests>
    {
        readonly PackageFile artifactsFile;

		public PackageFileTests() : base()
		{
            artifactsFile = Provider.GetRequiredService<PackageFile>();
        }

        [Fact]
        public void should_extract_to_DestinationDirectoryName()
        {
            artifactsFile.Extract();
            Assert.True(Directory.Exists(artifactsFile.ExtractedTo));
        }

        [Fact]
		public void should_set_full_control_of_extracted_folder()
		{
            artifactsFile.Extract();
            Assert.True(new DirectoryInfo(artifactsFile.ExtractedTo).UnixFileMode.HasFlag(UnixFileMode.UserExecute));
        }

        [Fact]
        public void hash_should_validate()
        {
            const string expectedHash = "3be95f70f2a217ac0278cd62d72fdf12732f9a69e92db77f2a70f5f3e6996e2a";

            var result = artifactsFile.IsValidHash(expectedHash);
            Assert.True(result);
        }

        public override void ConfigureServices()
        {
            var file = Test.DataFile.Get(PackageFile.FileName);
            Services.AddSingleton<PackageFile>(new PackageFile(file.FullName, Mock.Logger<PackageFile>()));
        }
    }
}

