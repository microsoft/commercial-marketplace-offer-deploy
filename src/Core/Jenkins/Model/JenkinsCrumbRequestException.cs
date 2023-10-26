using System;
namespace Modm.Jenkins.Model
{
    class JenkinsCrumbRequestException : Exception
    {
        public JenkinsCrumbRequestException(string message = null, Exception innerException = null) : base(message, innerException)
        {
        }

        public JenkinsCrumbRequestException(string message = null) : base(message)
        {
        }
    }
}

