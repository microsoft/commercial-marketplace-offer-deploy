package eventsfiltering

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

// filters resources
type EventGridEventFilter interface {
	// filters the event grid events by matching ANY of the filter tags
	// matchAny: at least 1 tag must match for the resource that the event grid event is for
	Filter(ctx context.Context, matchAny d.LookupTags, events []*eventgrid.Event) eg.EventGridEventResources
}

// creates a new tag filter
//
//	tags: the tags that will be used to return results (not filter on)
func NewTagsFilter(includeKeys []string, credential azcore.TokenCredential) EventGridEventFilter {
	mapper := newEventGridEventMapper(credential)
	return &tagsFilter{
		includeKeys: includeKeys,
		mapper:      mapper,
	}
}

// event grid events tag filter implementation
type tagsFilter struct {
	//search for keys that will cause the filter results to include these tags by key
	includeKeys []string
	mapper      eventGridEventMapper
}

// Filter implements ResourceFilter
func (f *tagsFilter) Filter(ctx context.Context, matchAny d.LookupTags, events []*eventgrid.Event) eg.EventGridEventResources {
	mappedItems := f.mapper.Map(ctx, events)
	items := []*eg.EventGridEventResource{}

	for _, item := range mappedItems {
		if matches, ok := f.match(matchAny, item.Resource.Tags); ok {
			item.Tags = matches
			items = append(items, item)
		}
	}

	return items
}

func (f *tagsFilter) match(matchAny d.LookupTags, resourceTags map[string]*string) (d.LookupTags, bool) {
	matches := d.LookupTags{}
	match := false

	// populate all tags that match searchFor
	for _, key := range f.includeKeys {
		if resourceTags[key] != nil {
			matches[d.LookupTagKey(key)] = resourceTags[key]
		}
	}

	for key := range matchAny {
		strKey := string(key)
		if resourceTags[strKey] != nil && strings.EqualFold(*matchAny[key], *resourceTags[strKey]) {
			match = true
			break
		}
	}

	return matches, match
}
