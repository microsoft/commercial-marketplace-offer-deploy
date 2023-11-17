using System.Security.Cryptography;
using Microsoft.Extensions.DependencyInjection;
using Modm.Packaging;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class PackageFileTests : AbstractTest<PackageFileTests>
    {
        readonly PackageFile file;

		public PackageFileTests() : base()
		{
            file = Provider.GetRequiredService<PackageFile>();
        }

        [Fact]
        public void should_extract_to_DestinationDirectoryName()
        {
            file.Extract();
            Assert.True(Directory.Exists(file.ExtractedTo));
        }

        [Fact]
		public void should_set_full_control_of_extracted_folder()
		{
            file.Extract();
            Assert.True(new DirectoryInfo(file.ExtractedTo).UnixFileMode.HasFlag(UnixFileMode.UserExecute));
        }

        [Fact]
        public void hash_should_validate()
        {
            const string expectedHash = "8016f746d03de6312283396c6e0f95504dcd14d58162f0f14bea28bf96c09663";
            Assert.True(file.IsValidHash(expectedHash));
        }

        protected override void ConfigureServices()
        {
            var file = Test.DataFile.Get(PackageFile.FileName);
            Services.AddSingleton<PackageFile>(new PackageFile(file.FullName, Mock.Logger<PackageFile>()));
        }
    }
}

