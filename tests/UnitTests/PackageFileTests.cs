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
            const string expectedHash = "692c18ed56f41ce23ac4296f482c66c3dee8b2b7d440ce2f4974d5a0adf63301";
            Assert.True(file.IsValidHash(expectedHash));
        }

        protected override void ConfigureServices()
        {
            var file = Test.DataFile.Get(PackageFile.FileName);
            Services.AddSingleton<PackageFile>(new PackageFile(file.FullName, Mock.Logger<PackageFile>()));
        }
    }
}

