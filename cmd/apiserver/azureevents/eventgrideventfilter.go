package azureevents

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

// filters resources
type EventGridEventFilter interface {
	// filters the event grid events by matching ANY of the filter tags
	// matchAny: at least 1 tag must match for the resource that the event grid event is for
	Filter(ctx context.Context, matchAny d.LookupTags, events []*eventgrid.Event) []*ResourceEventSubject
}

// creates a new tag filter
//
//	tags: the tags that will be used to return results (not filter on)
func NewTagsFilter(includeKeys []string, provider ResourceEventSubjectFactory) EventGridEventFilter {
	return &tagsFilter{
		includeKeys: includeKeys,
		factory:     provider,
	}
}

// event grid events tag filter implementation
type tagsFilter struct {
	//search for keys that will cause the filter results to include these tags by key
	includeKeys []string
	factory     ResourceEventSubjectFactory
}

// Filter implements ResourceFilter
func (f *tagsFilter) Filter(ctx context.Context, matchAny d.LookupTags, events []*eventgrid.Event) []*ResourceEventSubject {
	mappedItems := f.factory.Create(ctx, events)
	items := []*ResourceEventSubject{}

	for _, item := range mappedItems {
		if matches, ok := f.match(matchAny, item.Tags()); ok {
			item.SetLookupTags(matches)
			items = append(items, item)
		}

		if item.IsParentDeployment() {
			items = append(items, item)
		}
	}

	log.Debugf("factory received %d EventGridEvents, filtered to %d", len(events), len(items))

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
