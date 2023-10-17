namespace Modm.Packaging
{
    public interface IPackageDownloader
    {
        Task<PackageFile> DownloadAsync(PackageUri uri);
    }
}