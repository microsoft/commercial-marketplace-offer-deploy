// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
namespace Modm.Azure.Model
{
	public class UserDataResult
	{
		public UserData UserData { get; }
		public bool IsValid { get; }

        public UserDataResult(UserData userData)
		{
			UserData = userData;
			IsValid = userData != null && UserData.IsValid();
		}
	}
}

