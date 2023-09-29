namespace Modm.Artifacts
{
    public interface IArtifactsDownloader
    {
        Task<ArtifactsFile> DownloadAsync(ArtifactsUri uri);
    }
}