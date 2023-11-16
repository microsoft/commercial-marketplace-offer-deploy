using System;
using Modm.Deployments;
using Xunit.Abstractions;
using NSubstitute;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging.Abstractions;
using Modm.Extensions;

namespace Modm.Tests.UnitTests
{
	public class DeploymentRecordUnitTests
	{
        private readonly IConfiguration configuration;
        private readonly DeploymentFile deploymentFile;
        private readonly ITestOutputHelper output;
        private string tempPath;

        public DeploymentRecordUnitTests(ITestOutputHelper output)
		{
			this.output = output;
            this.configuration = Substitute.For<IConfiguration>();
            this.tempPath = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());

            this.configuration.GetHomeDirectory().Returns(tempPath);

            Directory.CreateDirectory(tempPath);

            this.deploymentFile = new DeploymentFile(this.configuration, new NullLogger<DeploymentFile>());
        }

        [Fact]
        public void ConfigurationSubstitute_ShouldReturnExpectedHomeDirectory()
        {
            // Arrange
            var expectedPath = this.tempPath;
            var configuration = Substitute.For<IConfiguration>();
            configuration.GetHomeDirectory().Returns(expectedPath);

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

            configuration.Received().GetHomeDirectory();

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

        public void Dispose()
        {
            if (Directory.Exists(this.tempPath))
            {
                Directory.Delete(this.tempPath, true);
            }
        }
    }
}

