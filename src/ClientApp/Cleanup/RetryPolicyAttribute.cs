using System;

namespace ClientApp.Cleanup
{
    [AttributeUsage(AttributeTargets.Class)]
    public class RetryPolicyAttribute : Attribute
    {
        public int RetryCount { get; set; } = 3;
        public int SleepDuration { get; set; } = 500;


        public Func<int, TimeSpan> GetSleepDurationProvider()
        {
            return new Func<int, TimeSpan>(i => TimeSpan.FromMilliseconds(i * this.SleepDuration));
        }
    }
}

