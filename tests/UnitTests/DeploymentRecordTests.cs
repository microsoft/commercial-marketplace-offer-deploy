using Modm.Deployments;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging.Abstractions;
using Modm.Extensions;
using Modm.Configuration;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class DeploymentRecordTests : IDisposable
	{
        private readonly IConfiguration configuration;
        private readonly DeploymentFile deploymentFile;
        private readonly DisposableDirectory<DeploymentRecordTests> tempDir;

        public DeploymentRecordTests()
		{
            this.tempDir = Test.Directory<DeploymentRecordTests>();

            this.configuration = new ConfigurationBuilder()
            .AddInMemoryCollection(new Dictionary<string, string?> {
                { EnvironmentVariable.Names.HomeDirectory, this.tempDir.FullName }
            }).Build();

            this.deploymentFile = new DeploymentFile(this.configuration, new NullLogger<DeploymentFile>());
        }

        [Fact]
        public void ConfigurationSubstitute_ShouldReturnExpectedHomeDirectory()
        {
            var actualPath = configuration.GetHomeDirectory();
            Assert.Equal(tempDir.FullName, actualPath);
        }

        [Fact]
        public async Task Write_DeploymentRecord_WritesToFile()
        {
            var initialDeployment = new Deployment { Id = 0, Timestamp = DateTimeOffset.UtcNow, Status = "Test" };
            var updatedDeployment = new Deployment { Id = 1, Timestamp = DateTimeOffset.UtcNow.AddSeconds(5), Status = "Test2" };

            var deploymentRecord = new DeploymentRecord(updatedDeployment);

            var auditRecord = new AuditRecord();
            auditRecord.AdditionalData.Add("initial", initialDeployment);
            deploymentRecord.AuditRecords.Add(auditRecord);

            var auditRecord2 = new AuditRecord();
            auditRecord2.AdditionalData.Add("updated", updatedDeployment);
            deploymentRecord.AuditRecords.Add(auditRecord2);

            await this.deploymentFile.Write(deploymentRecord, default);

            var filePath = Path.Combine(tempDir.FullName, DeploymentFile.FileName);
            Assert.True(File.Exists(filePath));
        }

        [Fact]
        public async Task Read_DeploymentRecord_ReadsFromFile()
        { 
            var initialDeployment = new Deployment { Id = 0, Timestamp = DateTimeOffset.UtcNow, Status = "Test" };
            var updatedDeployment = new Deployment { Id = 1, Timestamp = DateTimeOffset.UtcNow.AddSeconds(5), Status = "Test2" };

            var deploymentRecord = new DeploymentRecord(updatedDeployment);

            var auditRecord = new AuditRecord();
            auditRecord.AdditionalData.Add("initial", initialDeployment);
            deploymentRecord.AuditRecords.Add(auditRecord);

            var auditRecord2 = new AuditRecord();
            auditRecord2.AdditionalData.Add("updated", updatedDeployment);
            deploymentRecord.AuditRecords.Add(auditRecord2);

            await this.deploymentFile.Write(deploymentRecord, default);

            var readRecord = await deploymentFile.Read(default);

            Assert.Equal(1, readRecord.Deployment.Id);
            Assert.Equal("Test2", readRecord.Deployment.Status);
        }

        public void Dispose()
        {
            tempDir.Dispose();
        }
    }
}

