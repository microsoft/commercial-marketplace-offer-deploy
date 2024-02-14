// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using MediatR;

namespace Modm.Engine.Pipelines
{
    public abstract class Pipeline<TRequest, TResponse> : IPipeline<TRequest, TResponse>
        where TRequest : IRequest<TResponse>
    {
        protected readonly IMediator mediator;

        public Pipeline(IMediator mediator)
        {
            this.mediator = mediator;
        }

        /// <summary>
        /// Execute the pipeline with the specific request payload to get the response
        /// </summary>
        /// <param name="request"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        public virtual Task<TResponse> Execute(TRequest request, CancellationToken cancellationToken = default)
        {
            return mediator.Send(request, cancellationToken);
        }
    }
}

