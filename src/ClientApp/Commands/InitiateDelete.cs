// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using MediatR;
using Modm.Azure;
using ClientApp.Backend;

namespace ClientApp.Commands
{
    public class InitiateDelete : IRequest
    {
        private readonly string resourceGroupName;

        public InitiateDelete(string resourceGroupName)
        {
            this.resourceGroupName = resourceGroupName;
        }

        public string ResourceGroupName
        {
            get { return this.resourceGroupName; }
        }
    }
}

