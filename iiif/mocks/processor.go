package mocks

import iiif "github.com/t11e/picaxe/iiif"
import io "io"
import mock "github.com/stretchr/testify/mock"
import resources "github.com/t11e/picaxe/resources"

// Processor is an autogenerated mock type for the Processor type
type Processor struct {
	mock.Mock
}

// Process provides a mock function with given fields: req, resolver, w, result
func (_m *Processor) Process(req iiif.Request, resolver resources.Resolver, w io.Writer, result *iiif.Result) error {
	ret := _m.Called(req, resolver, w, result)

	var r0 error
	if rf, ok := ret.Get(0).(func(iiif.Request, resources.Resolver, io.Writer, *iiif.Result) error); ok {
		r0 = rf(req, resolver, w, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

var _ iiif.Processor = (*Processor)(nil)
