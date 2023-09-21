using System;
using FluentValidation;
using MediatR;

namespace Modm.Engine.Behaviors
{
    public class ValidationBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    {
        private readonly IValidator<TRequest> validator;

        public ValidationBehavior(IValidator<TRequest> validator)
        {
            this.validator = validator;
        }

        public Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
        {
            validator.ValidateAndThrow(request);
            return next();
        }
    }
}

