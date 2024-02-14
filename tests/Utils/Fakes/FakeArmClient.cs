// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using Azure.ResourceManager;
using NSubstitute;

namespace Modm.Tests.Utils.Fakes
{
    public class FakeArmClient : ArmClient
    {
        public FakeArmClient() : base()
        {
        }

        public static ArmClient New()
        {
            var instance = new FakeArmClient();
            return instance;
        }
    }
}

