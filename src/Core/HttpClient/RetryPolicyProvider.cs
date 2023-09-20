using System;
using Microsoft.Extensions.Logging;
using Polly;
using Polly.Retry;

namespace Modm.HttpClient
{
	public class RetryPolicyProvider
	{
        public static AsyncRetryPolicy GetRetryPolicy(ILogger logger)
        {
            return Policy
                .Handle<Exception>() // Specify the exceptions to handle for retries
                .WaitAndRetryAsync(
                    retryCount: 3, // Number of retries
                    sleepDurationProvider: attempt => TimeSpan.FromSeconds(Math.Pow(2, attempt)), // Incremental backoff (2^n seconds)
                    onRetry: (exception, calculatedWaitDuration, context) =>
                    {
                        // Log the retry attempt using ILogger
                        logger.LogWarning($"Retrying after {calculatedWaitDuration.TotalSeconds} seconds due to {exception.GetType().Name}: {exception.Message}");
                    });
        }
    }
}

