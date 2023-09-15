using System;
namespace Modm.ServiceHost
{
	public static class IServiceCollectionExtensions
	{
		public static IServiceCollection AddSingletonHostedService<T>(this IServiceCollection services) where T : class, IHostedService
		{
			services.AddSingleton<T>().AddHostedService(p => p.GetRequiredService<T>());
			return services;
        }
	}
}

