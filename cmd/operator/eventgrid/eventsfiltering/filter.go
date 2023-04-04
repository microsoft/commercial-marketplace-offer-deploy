package eventsfiltering

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

// filters resources
type EventGridEventFilter interface {
	// filters the event grid events by matching ANY of the filter tags
	// matchAny: at least 1 tag must match for the resource that the event grid event is for
	Filter(ctx context.Context, matchAny FilterTags, events []*eventgrid.Event) FilterResult
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
func (f *tagsFilter) Filter(ctx context.Context, matchAny FilterTags, events []*eventgrid.Event) FilterResult {
	mappedItems := f.mapper.Map(ctx, events)
	items := []FilterResultItem{}

	for _, item := range mappedItems {
		if matches, ok := f.match(matchAny, item.resource.Tags); ok {
			items = append(items, FilterResultItem{
				EventGridEvent: *item.event,
				Resource:       *item.resource,
				MatchedTags:    matches,
			})
		}
	}

	return FilterResult{
		Items: items,
	}
}

func (f *tagsFilter) match(matchAny FilterTags, resourceTags map[string]*string) (map[string]string, bool) {
	matches := map[string]string{}
	match := false

	// populate all tags that match searchFor
	for _, key := range f.includeKeys {
		if resourceTags[key] != nil {
			matches[key] = *resourceTags[key]
		}
	}

	for key, _ := range matchAny {
		if resourceTags[key] != nil && strings.EqualFold(*matchAny[key], *resourceTags[key]) {
			match = true
			break
		}
	}

	return matches, match
}
