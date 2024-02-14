// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm.Security
{
    /// <summary>
    /// Admin credentials that are set from Installer step in the createUiDefinition.json,
    /// used to login to the frontend dashboard
    /// </summary>
    public class AdminCredentials
    {
        public string Username { get; set; }
        public string Password { get; set; }
    }
}