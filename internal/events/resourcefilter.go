package events

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

// filters resources
type ResourceFilter interface {
	Filter(resource []*armresources.GenericResource) []*armresources.GenericResource
}

type ResourceTagFilter struct {
	Tags map[string]string
}

// Filter implements ResourceFilter
func (*ResourceTagFilter) Filter(resource []*armresources.GenericResource) []*armresources.GenericResource {
	panic("unimplemented")
}

func NewResourceTagFilter(tags map[string]string) ResourceFilter {
	return &ResourceTagFilter{
		Tags: tags,
	}
}
