using System;
using Modm.Azure;

namespace Modm.Tests.UnitTests
{
	public class ArmFunctionsTests
	{
		[Fact]
		public void UniqueStringShouldEqualValueCreatedFromArmTemplate()
		{
			const string resourceIdFromArmOutput = "/subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/rg-modm201-20230913120256/providers/Microsoft.Compute/virtualMachines/bobjacmodm201";
            const string valueFromArmTemplate = "yvjxno3je5on2";

			Assert.Equal(valueFromArmTemplate, ArmFunctions.UniqueString(resourceIdFromArmOutput));

        }
	}
}

