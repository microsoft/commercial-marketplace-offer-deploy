namespace Modm.Deployments
{
    public interface IDeploymentParametersFile
    {
        Task Write(IDictionary<string, object> parameters);
    }
}