namespace Modm.Deployments
{
    public interface IDeploymentParametersFile
    {
        string FullPath { get; }
        Task Write(IDictionary<string, object> parameters);
    }
}