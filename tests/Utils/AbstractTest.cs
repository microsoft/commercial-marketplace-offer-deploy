// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using Microsoft.Extensions.DependencyInjection;

namespace Modm.Tests.Utils
{
    public abstract partial class AbstractTest<TTest> : IDisposable
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

        /// <summary>
        /// Loads a service type up to then perform assertions with it using the
        /// instance that comes from the Provider built into the test
        /// </summary>
        /// <typeparam name="TService"></typeparam>
        /// <param name="action"></param>
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
    }
}