using Azure.Core;
using Azure.ResourceManager.Resources;
using MediatR;

namespace ClientApp.Cleanup
{
    public interface IDeleteResourceRequest : IRequest<DeleteResourceResult>
	{
        ResourceGroupResource ResourceGroup { get; }
        ResourceIdentifier ResourceId { get; }
	}
}