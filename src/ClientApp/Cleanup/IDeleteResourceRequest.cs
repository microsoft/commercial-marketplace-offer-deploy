using Azure.Core;
using MediatR;

namespace ClientApp.Cleanup
{
    public interface IDeleteResourceRequest : IRequest<DeleteResourceResult>
	{
        ResourceIdentifier ResourceId { get; }
	}
}