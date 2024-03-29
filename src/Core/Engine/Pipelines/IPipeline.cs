﻿using System;
using MediatR;

namespace Modm.Engine.Pipelines
{
	public interface IPipeline<TRequest, TResponse>
		where TRequest : IRequest<TResponse>
	{
		Task<TResponse> Execute(TRequest request, CancellationToken cancellationToken = default);
	}
}

