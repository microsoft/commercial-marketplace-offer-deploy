// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using MediatR;

namespace Modm.ServiceHost.Notifications
{
    /// <summary>
    /// internal notification that the user data is available on the vm
    /// </summary>
    public class UserDataAvailable : INotification
    {
        public required string UserData { get; set; }
    }
}

