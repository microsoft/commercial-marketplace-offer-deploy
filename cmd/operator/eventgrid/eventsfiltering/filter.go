package eventsfiltering

import (
	"context"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

// filters resources
type EventGridEventFilter interface {
	Filter(ctx context.Context, events []*eventgrid.Event) []*eventgrid.Event
}

// creates a new tag filter
func NewTagsFilter(tags FilterTags, credential azcore.TokenCredential) EventGridEventFilter {
	mapper := newEventGridEventMapper(credential)
	return &tagsFilter{
		tags:   tags,
		mapper: mapper,
	}
}

// event grid events tag filter implementation
type tagsFilter struct {
	tags   FilterTags
	mapper eventGridEventMapper
}

// Filter implements ResourceFilter
func (f *tagsFilter) Filter(ctx context.Context, events []*eventgrid.Event) []*eventgrid.Event {
	mappedItems := f.mapper.Map(ctx, events)
	filteredEvents := []*eventgrid.Event{}

	for _, item := range mappedItems {
		if f.match(item.resource.Tags) {
			filteredEvents = append(filteredEvents, item.event)
		}
	}
	return filteredEvents
}

func (f *tagsFilter) match(resourceTags map[string]*string) bool {
	return reflect.DeepEqual(f.tags, resourceTags)
}
