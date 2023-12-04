using System;
using ClientApp.Commands;
using MediatR;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Modm.ClientApp.Controllers;

namespace ClientApp.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class DeleteController : ControllerBase
    {
        private readonly IMediator mediator;
        private readonly ILogger<DeleteController> logger;

        public DeleteController(IMediator mediator, ILogger<DeleteController> logger)
		{
            this.mediator = mediator;
            this.logger = logger;
        }

        [HttpPost]
        [Route("resources/{resourceGroupName}/deletemodmresources")]
        public async Task<IActionResult> DeleteResourcesWithTagAsync([FromRoute] string resourceGroupName)
        {
            this.logger.LogInformation("Delete request received");
            var initiateDelete = new InitiateDelete(resourceGroupName);
            await this.mediator.Send(initiateDelete);

            return Ok("Successfully submitted a delete");
        }

    }
}

