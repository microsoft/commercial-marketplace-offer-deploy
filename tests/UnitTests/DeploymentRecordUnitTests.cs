using System;
using Modm.Deployments;
using Xunit.Abstractions;
using NSubstitute;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging.Abstractions;
using Modm.Extensions;
using Modm.Configuration;

namespace Modm.Tests.UnitTests
{
	public class DeploymentRecordUnitTests : IDisposable
	{
        private readonly IConfiguration configuration;
        private readonly DeploymentFile deploymentFile;
        private readonly ITestOutputHelper output;
        private string tempPath;

        public DeploymentRecordUnitTests(ITestOutputHelper output)
		{
			this.output = output;
            
            this.tempPath = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());

            var inMemorySettings = new Dictionary<string, string> {
                {EnvironmentVariable.Names.HomeDirectory, this.tempPath}
            };

            this.configuration = new ConfigurationBuilder()
                .AddInMemoryCollection(inMemorySettings)
                .Build();

            Directory.CreateDirectory(tempPath);

            this.deploymentFile = new DeploymentFile(this.configuration, new NullLogger<DeploymentFile>());
        }

        [Fact]
        public void ConfigurationSubstitute_ShouldReturnExpectedHomeDirectory()
        {
            // Arrange
            var expectedPath = this.tempPath;

            // Act
            var actualPath = configuration.GetHomeDirectory();

            // Assert
            Assert.Equal(expectedPath, actualPath);
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

            var filePath = Path.Combine(tempPath, DeploymentFile.FileName);
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

        [Fact]
        public async Task Update_DeploymentRecord_ReadsWritesFile()
        {
            var initialDeployment = new Deployment { Id = 1, Timestamp = DateTimeOffset.UtcNow, Status = "Initial" };
            var deploymentRecord = new DeploymentRecord(initialDeployment);

            var auditRecord = new AuditRecord();
            auditRecord.AdditionalData.Add("initial", initialDeployment);
            deploymentRecord.AuditRecords.Add(auditRecord);

            await this.deploymentFile.Write(deploymentRecord, default);

            var readRecord = await deploymentFile.Read(default);

            Assert.Equal(1, readRecord.Deployment.Id);
            Assert.Equal("Initial", readRecord.Deployment.Status);
            Assert.Single(readRecord.AuditRecords);

            var readDeployment = readRecord.Deployment;
            readDeployment.Status = "Updated";

            var updatedAuditRecord = new AuditRecord();
            updatedAuditRecord.AdditionalData.Add("updated", readDeployment);

            readRecord.AuditRecords.Add(updatedAuditRecord);
            await this.deploymentFile.Write(readRecord, default);


            var readUpdatedRecord = await deploymentFile.Read(default);
            var readUpdatedDeployment = readUpdatedRecord.Deployment;

            Assert.Equal(1, readUpdatedDeployment.Id);
            Assert.Equal("Updated", readUpdatedDeployment.Status);
            Assert.Equal(2, readUpdatedRecord.AuditRecords.Count);
        }

        public void Dispose()
        {
            if (Directory.Exists(this.tempPath))
            {
                Directory.Delete(this.tempPath, true);
            }
        }
    }
}

