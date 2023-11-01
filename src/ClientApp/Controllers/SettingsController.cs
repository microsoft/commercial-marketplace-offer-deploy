using System;
using Microsoft.AspNetCore.Mvc;

namespace Modm.ClientApp.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class SettingsController : ControllerBase
    {
        private readonly IConfiguration configuration;

        public SettingsController(IConfiguration configuration)
        {
            this.configuration = configuration;
        }

        [HttpGet()]
        public object? Get([FromQuery] string key)
        {
            if (string.IsNullOrEmpty(key))
            {
                return default;
            }
            return configuration.GetValue<object>(key.Replace(".", ":"));
        }
    }
}