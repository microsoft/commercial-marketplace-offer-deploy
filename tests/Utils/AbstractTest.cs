using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Modm.Packaging;
using Modm.Configuration;
using Modm.Deployments;
using Modm.Jenkins.Client;
using NSubstitute;

namespace Modm.Tests.Utils
{
    public abstract class AbstractTest<TTest> : IDisposable
    {
        private HashSet<IDisposable> disposables { get; } = new();
        private readonly ServiceProvider provider;

        protected ServiceProvider Provider => provider;
        protected ServiceCollection Services { get; } = new();

        protected MockConfigurator Mock => new(Services, disposables);

        public AbstractTest()
        {
            ConfigureServices();
            provider = Services.BuildServiceProvider();
        }

        protected abstract void ConfigureServices();

        protected void ConfigureMocks(Action<MockConfigurator> configure)
        {
            configure(Mock);
        }


        protected void With<TService>(Action<TService> action) where TService : notnull
        {
            var instance = Provider.GetRequiredService<TService>();
            action(instance);
        }

        public virtual void Dispose()
        {
            provider?.Dispose();
            foreach (var disposable in disposables)
            {
                disposable.Dispose();
            }
        }

        public class MockConfigurator
        {
            private readonly ServiceCollection services;
            private readonly HashSet<IDisposable> disposables;

            public MockConfigurator(ServiceCollection services, HashSet<IDisposable> disposables)
            {
                this.services = services;
                this.disposables = disposables;
            }

            public T Create<T>(Action<T>? configure = default) where T: class
            {
                var instance = Substitute.For<T>();
                configure?.Invoke(instance);
                return instance;
            }

            public ILogger<T> Logger<T>()
            {
                var instance = Substitute.For<ILogger<T>>();
                services.AddSingleton(instance);

                return instance;
            }

            public IJenkinsClient JenkinsClient(Action<IJenkinsClient>? configure = default)
            {
                var instance = Substitute.For<IJenkinsClient>();
                configure?.Invoke(instance);

                services.AddSingleton<IJenkinsClient>(instance);

                return instance;
            }

            public JenkinsClientFactory JenkinsClientFactory(Action<JenkinsClientFactory>? configure = default)
            {
                var instance = Substitute.For<JenkinsClientFactory>();
                configure?.Invoke(instance);

                services.AddSingleton<JenkinsClientFactory>(instance);

                return instance;
            }

            public IConfiguration Configuration()
            {
                var dir = Test.Directory<TTest>();
                disposables.Add(dir);

                var configuration = Substitute.For<IConfiguration>();
                configuration.GetValue<string>(EnvironmentVariable.Names.HomeDirectory).Returns(dir.FullName);

                services.AddSingleton<IConfiguration>(configuration);

                return configuration;
            }

            public IPackageDownloader PackageDownloader()
            {
                var file = Test.DataFile.Get(PackageFile.FileName);
                var dir = Test.Directory<TTest>();
                disposables.Add(dir);

                // copy our file to the temp dir
                var filePath = Path.Combine(dir.FullName, file.Name);
                File.Copy(file.FullName, filePath, true);

                var instance = Substitute.For<IPackageDownloader>();
                instance.DownloadAsync(Arg.Any<PackageUri>())
                    .Returns(new PackageFile(filePath, Logger<PackageFile>()));

                services.AddSingleton(instance);

                return instance;
            }

            public IDeploymentRepository DeploymentRepository()
            {
                var instance = Substitute.For<IDeploymentRepository>();
                instance.Get().Returns(new Deployment());

                services.AddSingleton(instance);

                return instance;
            }
        }
    }
}

