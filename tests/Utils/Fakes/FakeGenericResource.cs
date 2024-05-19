using Azure.ResourceManager.Resources;
using NSubstitute;

namespace Modm.Tests.Utils.Fakes
{
    public class FakeGenericResource : GenericResource
    {
        public static GenericResource New(Action<FakeGenericResource>? configure = default)
        {
            var instance = Substitute.ForPartsOf<FakeGenericResource>();
            configure?.Invoke(instance);
            return instance;
        }
    }
}

