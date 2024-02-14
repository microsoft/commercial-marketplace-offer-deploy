// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Modm.Deployments;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging.Abstractions;
using Modm.Extensions;
using Modm.Configuration;
using Modm.Tests.Utils;

namespace Modm.Tests.UnitTests
{
    public class DeploymentFileTests : IDisposable
	{
        private readonly IConfiguration configuration;
        private readonly DeploymentFile deploymentFile;
        private readonly DisposableDirectory<DeploymentFileTests> tempDir;

        public DeploymentFileTests()
		{
            this.tempDir = Test.Directory<DeploymentFileTests>();

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

           // var deploymentRecord = new DeploymentRecord(updatedDeployment);

            //var auditRecord = new AuditRecord();
            //auditRecord.AdditionalData.Add("initial", initialDeployment);
            //deploymentRecord.AuditRecords.Add(auditRecord);

            //var auditRecord2 = new AuditRecord();
            //auditRecord2.AdditionalData.Add("updated", updatedDeployment);
            //deploymentRecord.AuditRecords.Add(auditRecord2);

            await this.deploymentFile.WriteAsync(initialDeployment, default);

            var filePath = Path.Combine(tempDir.FullName, "deployment.json");
            Assert.True(File.Exists(filePath));
        }

        [Fact]
        public async Task Read_DeploymentRecord_ReadsFromFile()
        { 
            var initialDeployment = new Deployment { Id = 0, Timestamp = DateTimeOffset.UtcNow, Status = "Test" };
            //var updatedDeployment = new Deployment { Id = 1, Timestamp = DateTimeOffset.UtcNow.AddSeconds(5), Status = "Test2" };

            //var deploymentRecord = new DeploymentRecord(updatedDeployment);

            //var auditRecord = new AuditRecord();
            //auditRecord.AdditionalData.Add("initial", initialDeployment);
            //deploymentRecord.AuditRecords.Add(auditRecord);

            //var auditRecord2 = new AuditRecord();
            //auditRecord2.AdditionalData.Add("updated", updatedDeployment);
            //deploymentRecord.AuditRecords.Add(auditRecord2);

            await this.deploymentFile.WriteAsync(initialDeployment, default);

            var readDeployment = await deploymentFile.ReadAsync(default);

            Assert.Equal(0, readDeployment.Id);
            Assert.Equal("Test", readDeployment.Status);
        }

        [Fact]
        public async Task Update_DeploymentRecord_ReadsWritesFile()
        {
            var initialDeployment = new Deployment { Id = 1, Timestamp = DateTimeOffset.UtcNow, Status = "Initial" };
            //var deploymentRecord = new DeploymentRecord(initialDeployment);

            //var auditRecord = new AuditRecord();
            //auditRecord.AdditionalData.Add("initial", initialDeployment);
            //deploymentRecord.AuditRecords.Add(auditRecord);

            await this.deploymentFile.WriteAsync(initialDeployment, default);

            var readDeployment = await deploymentFile.ReadAsync(default);

            Assert.Equal(1, readDeployment.Id);
            Assert.Equal("Initial", readDeployment.Status);
            //Assert.Single(readRecord.AuditRecords);

            //var readDeployment = readRecord.Deployment;
            readDeployment.Status = "Updated";

            //var updatedAuditRecord = new AuditRecord();
            //updatedAuditRecord.AdditionalData.Add("updated", readDeployment);

            //readRecord.AuditRecords.Add(updatedAuditRecord);
            await this.deploymentFile.WriteAsync(readDeployment, default);

            var readUpdatedDeployment = await deploymentFile.ReadAsync(default);

            Assert.Equal(1, readUpdatedDeployment.Id);
            Assert.Equal("Updated", readUpdatedDeployment.Status);
        }

        public void Dispose()
        {
            tempDir.Dispose();
        }
    }
}

