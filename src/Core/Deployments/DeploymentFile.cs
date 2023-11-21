using System.Text.Json;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Modm.Extensions;

namespace Modm.Deployments
{
    using Microsoft.Extensions.Configuration;
    using Microsoft.Extensions.Logging;

    
    public class DeploymentFile : JsonFile<Deployment>
    {
        public override string FileName => "deployment.json";

        public DeploymentFile(IConfiguration configuration, ILogger<DeploymentFile> logger)
            : base(configuration, logger)
        {
        }
    }
    

    //   public class DeploymentFile
    //{
    //       public const string FileName = "deployment.json";

    //       public DateTimeOffset Timestamp
    //       {
    //           get
    //           {
    //               return new FileInfo(GetDeploymentFilePath()).LastWriteTimeUtc;
    //           }
    //       }

    //       private static readonly JsonSerializerOptions serializerOptions = new JsonSerializerOptions
    //       {
    //           PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
    //           WriteIndented = true
    //       };

    //       private readonly IConfiguration configuration;
    //       private readonly ILogger<DeploymentFile> logger;

    //       public DeploymentFile(IConfiguration configuration, ILogger<DeploymentFile> logger)
    //	{
    //           this.configuration = configuration;
    //           this.logger = logger;
    //       }

    //       public async Task<DeploymentRecord> Read(CancellationToken cancellationToken = default)
    //       {
    //           var path = GetDeploymentFilePath();

    //           if (!File.Exists(path))
    //           {
    //               this.logger.LogError($"{path} was not found");
    //               var deployment = new Deployment
    //               {
    //                   Id = 0,
    //                   Timestamp = DateTimeOffset.UtcNow,
    //                   Status = DeploymentStatus.Undefined
    //               };

    //               var newRecord = new DeploymentRecord(deployment);
    //               var auditRecord = new AuditRecord();
    //               auditRecord.AdditionalData.Add("initialState", deployment);
    //               newRecord.AuditRecords.Add(auditRecord);
    //               return newRecord;
    //           }

    //           var json = await File.ReadAllTextAsync(GetDeploymentFilePath(), cancellationToken);
    //           var deploymentRecord = JsonSerializer.Deserialize<DeploymentRecord>(json, serializerOptions);

    //           return deploymentRecord;
    //       }

    //       public async Task Write(DeploymentRecord deploymentRecord, CancellationToken cancellationToken)
    //       {
    //           var json = JsonSerializer.Serialize(deploymentRecord, serializerOptions);
    //           this.logger.LogInformation($"Writing Deployment json - {json}");
    //           await File.WriteAllTextAsync(GetDeploymentFilePath(), json, cancellationToken);
    //           this.logger.LogInformation($"Wrote the deployment json to {GetDeploymentFilePath()}");
    //       }

    //       private string GetDeploymentFilePath()
    //       {
    //           return Path.GetFullPath(Path.Combine(configuration.GetHomeDirectory(), FileName));
    //       }
    //   }
}

