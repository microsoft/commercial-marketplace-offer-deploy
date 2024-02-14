// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm.Packaging
{
    public interface IPackageDownloader
    {
        Task<PackageFile> DownloadAsync(PackageUri uri);
    }
}