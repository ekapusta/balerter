package registry

import (
	"fmt"
	"github.com/balerter/balerter/internal/lua_formatter"
	"github.com/balerter/balerter/internal/modules"
)

func (r *Registry) Result() []modules.TestResult {
	results := make([]modules.TestResult, 0)

	// iterate over calls
	for _, call := range r.calls {
		// an assert for this call is not registered
		e, ok := r.getAssert(call)
		if !ok {
			continue
		}
		// all asserts for this call already processed
		if len(e.asserts) == 0 {
			continue
		}

		shouldBeCalled := e.asserts[0]
		e.asserts = e.asserts[1:]

		res := modules.TestResult{}
		if shouldBeCalled {
			res.Message = fmt.Sprintf("method '%s' with args %s was called", call.method, lua_formatter.ValuesToStringNoErr(call.args))
			res.Ok = true
		} else {
			res.Message = fmt.Sprintf("method '%s' with args %s was called, but should not", call.method, lua_formatter.ValuesToStringNoErr(call.args))
			res.Ok = false
		}
		results = append(results, res)
	}

	results = append(results, r.getAssertsOrphans()...)

	return results
}

func (r *Registry) getAssert(call call) (*assertEntry, bool) {
	e, ok := r.assertEntries[call.method]
	if !ok {
		return nil, false
	}

	for _, a := range call.args {
		key := lua_formatter.ValueToStringNoErr(a)
		e1, ok := e.entries[key]
		if !ok {
			return nil, false
		}
		e = e1
	}

	return e, true
}

func (r *Registry) getAssertsOrphans() []modules.TestResult {
	result := make([]modules.TestResult, 0)

	for method, e := range r.assertEntries {
		result = append(result, r.getAssertsOrphansChain(e, method)...)
	}

	return result
}

func (r *Registry) getAssertsOrphansChain(entry *assertEntry, method string, args ...string) []modules.TestResult {
	result := make([]modules.TestResult, 0)

	for _, v := range entry.asserts {
		res := modules.TestResult{}
		if v {
			res.Message = fmt.Sprintf("method '%s' with args %v was not called, but should", method, args)
			res.Ok = false
		} else {
			res.Message = fmt.Sprintf("method '%s' with args %v was not called", method, args)
			res.Ok = true
		}
		result = append(result, res)
	}

	for arg, e := range entry.entries {
		results := r.getAssertsOrphansChain(e, method, append(args, arg)...)
		result = append(result, results...)
	}

	return result
}
