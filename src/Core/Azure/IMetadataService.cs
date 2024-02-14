// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Modm.Azure.Model;

namespace Modm.Azure
{
    public interface IMetadataService
    {
        Task<InstanceMetadata> GetAsync();
        Task<string> GetFqdnAsync();
        Task<UserDataResult> TryGetUserData();
    }
}