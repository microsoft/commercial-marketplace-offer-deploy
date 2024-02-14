// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace ClientApp.Cleanup
{
    public interface IDeleteProcessor
	{
		Task DeleteResourcesAsync(string resourceGroup, CancellationToken cancellationToken = default);
	}
}