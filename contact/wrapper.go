package contact

import (
	"google.golang.org/api/googleapi"
	"google.golang.org/api/people/v1"
)

type PeopleApiWrapper struct {
	call *people.PeopleConnectionsListCall
}

func (wrapper *PeopleApiWrapper) PageToken(pageToken string) ServiceInterface {
	wrapper.call.PageToken(pageToken)
	return wrapper
}

func (wrapper *PeopleApiWrapper) PersonFields(personFields string) ServiceInterface {
	wrapper.call.PersonFields(personFields)
	return wrapper
}

func (wrapper *PeopleApiWrapper) Do(opts ...googleapi.CallOption) (*people.ListConnectionsResponse, error) {
	return wrapper.call.Do(opts...)
}
