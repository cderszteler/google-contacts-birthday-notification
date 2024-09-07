package contact

import (
	"google.golang.org/api/googleapi"
	"google.golang.org/api/people/v1"
)

type peopleApiWrapper struct {
	call *people.PeopleConnectionsListCall
}

func (wrapper *peopleApiWrapper) PageToken(pageToken string) serviceInterface {
	wrapper.call.PageToken(pageToken)
	return wrapper
}

func (wrapper *peopleApiWrapper) PersonFields(personFields string) serviceInterface {
	wrapper.call.PersonFields(personFields)
	return wrapper
}

func (wrapper *peopleApiWrapper) Do(opts ...googleapi.CallOption) (*people.ListConnectionsResponse, error) {
	return wrapper.call.Do(opts...)
}
