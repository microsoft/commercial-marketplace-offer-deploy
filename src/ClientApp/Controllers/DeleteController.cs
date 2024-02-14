// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
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
            this.logger.LogInformation("Dispatching initiateDelete");

            try
            {
                await this.mediator.Send(initiateDelete);
                return Ok(new { Message = "Delete operation successfully submitted." });
            }
            catch (Exception ex)
            {
                this.logger.LogError(ex, "Error initiating delete");
                return StatusCode(StatusCodes.Status500InternalServerError, new { Message = "Delete operation failed." });
            }
        }
    }
}

