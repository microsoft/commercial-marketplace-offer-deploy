using System;

namespace Modm.Deployments
{
	public class ArmParametersFile : IDeploymentParametersFile
    {
        public const string FileName = "parameters.json";
        private readonly string destinationDirectory;

        public ArmParametersFile(string destinationDirectory)
		{
            this.destinationDirectory = destinationDirectory;
		}

        public string FullPath => Path.Combine(destinationDirectory, FileName);

        public async Task Write(IDictionary<string, object> parameters)
        {
            string json = parameters.ToArmParametersJson();
            await File.WriteAllTextAsync(FullPath, json);
        }
    }
}

