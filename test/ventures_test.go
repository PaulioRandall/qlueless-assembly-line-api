package test

import (
	"bytes"
	"encoding/json"
	"testing"

	a "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/asserts"
	v "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

// _beginVenTest is run at the start of every test to setup the server and
// inject the test data.
func _beginVenTest() {
	venDBReset()
	venDBInjectVentures()
	startServer()
}

// _endVenTest should be deferred straight after _beginVenTest() is run to
// close resources at the end of every test.
func _endVenTest() {
	stopServer()
	venDBClose()
}

// ****************************************************************************
// (GET) /ventures
// ****************************************************************************

func TestGET_Ventures_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When all Ventures are requested
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON array of valid Ventures
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	v.AssertVentureSliceFromReader(t, res.Body)
}

// ****************************************************************************
// (GET) /ventures?wrap
// ****************************************************************************

func TestGET_Ventures_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When all Ventures are requested
		And the 'wrap' query parameter has been specified
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON array of valid Ventures wrapped with meta information
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures?wrap",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	v.AssertWrappedVentureSliceFromReader(t, res.Body)
}

// ****************************************************************************
// (GET) /ventures?id={id}
// ****************************************************************************

func TestGET_Venture_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a specific existing Venture is requested
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing a valid Venture
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures?id=1",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	v.AssertVentureFromReader(t, res.Body)
}

func TestGET_Venture_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a specific non-existent Venture is requested
		Then ensure the response code is 404
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures?id=999999",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 400, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	assertWrappedErrorBody(t, res.Body)
}

// ****************************************************************************
// (GET) /ventures?wrap&id={id}
// ****************************************************************************

func TestGET_Venture_3(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a specific existing Venture is requested
		And the 'wrap' query parameter has been specified
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing a valid Venture wrapped with meta information
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures?wrap&id=1",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	v.AssertWrappedVentureFromReader(t, res.Body)
}

// ****************************************************************************
// (POST) /ventures
// ****************************************************************************

func TestPOST_Venture_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a new valid Venture is POSTed
		Then ensure the response code is 201
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing the living input Venture with a new assigned ID
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.Venture{
		Description: "A new Venture",
		State:       "Not started",
		OrderIDs:    "1,2,3",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "POST",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 201, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	output := v.AssertVentureFromReader(t, res.Body)

	input.ID = output.ID
	input.IsDead = false
	v.AssertGenericVenture(t, output)
	v.AssertVentureModEquals(t, input, output)
}

func TestPOST_Venture_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a new invalid Venture is POSTed
		Then ensure the response code is 400
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.Venture{
		Description: "",
		State:       "",
		OrderIDs:    "invalid",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "POST",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 400, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	assertWrappedErrorBody(t, res.Body)
}

// ****************************************************************************
// (POST) /ventures?wrap
// ****************************************************************************

func TestPOST_Venture_3(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a new valid Venture is POSTed
		And the 'wrap' query parameter has been specified
		Then ensure the response code is 201
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing a WrappedReply
		And that the wrapped data is the living input Venture with a new assigned ID
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.Venture{
		Description: "A new Venture",
		State:       "Not started",
		OrderIDs:    "1,2,3",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures?wrap",
		Method: "POST",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 201, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	_, output := v.AssertWrappedVentureFromReader(t, res.Body)

	input.ID = output.ID
	input.IsDead = false
	v.AssertGenericVenture(t, output)
	v.AssertVentureModEquals(t, input, output)
}

// ****************************************************************************
// (PUT) /ventures
// ****************************************************************************

func TestPUT_Ventures_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When an existing Venture is modified and PUT to the server
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing the updated input Venture
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		IDs:   "1",
		Props: "description, state, order_ids, extra",
		Values: v.Venture{
			Description: "Black blizzard",
			State:       "In progress",
			OrderIDs:    "1,2,3",
			Extra:       "colour: black; power: 9000",
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	output := v.AssertVentureSliceFromReader(t, res.Body)
	require.Len(t, output, 1)

	input.Values.ID = "1"
	input.Values.IsDead = false
	v.AssertGenericVenture(t, output[0])
	v.AssertVentureModEquals(t, input.Values, output[0])
}

func TestPUT_Ventures_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When an non-existent Venture is PUT to the server
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing an empty Venture array
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		IDs:   "999999",
		Props: "description, state, order_ids, extra",
		Values: v.Venture{
			Description: "Black blizzard",
			State:       "In progress",
			OrderIDs:    "1,2,3",
			Extra:       "colour: black; power: 9000",
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	output := v.AssertVentureSliceFromReader(t, res.Body)
	require.Empty(t, output)
}

func TestPUT_Ventures_3(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a venture modification without IDs is PUT to the server
		Then ensure the response code is 400
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		Props: "description, state, order_ids, extra",
		Values: v.Venture{
			Description: "Black blizzard",
			State:       "In progress",
			OrderIDs:    "1,2,3",
			Extra:       "colour: black; power: 9000",
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 400, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	assertWrappedErrorBody(t, res.Body)
}

func TestPUT_Ventures_4(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When ventures updates are PUT to the server with invalid modifications
		Then ensure the response code is 400
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		IDs:    "1",
		Props:  "description, state, order_ids, extra",
		Values: v.Venture{},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 400, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)
	assertWrappedErrorBody(t, res.Body)
}

func TestPUT_Ventures_5(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When existing Ventures are modified as dead and PUT to the server
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing the updated input Venture
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		IDs:   "4,5",
		Props: "is_dead",
		Values: v.Venture{
			IsDead: true,
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	output := v.AssertVentureSliceFromReader(t, res.Body)
	require.Len(t, output, 2)

	assert.Equal(t, "4", output[0].ID)
	assert.True(t, output[0].IsDead)

	assert.Equal(t, "5", output[1].ID)
	assert.True(t, output[1].IsDead)
}

// ****************************************************************************
// (PUT) /ventures?wrap
// ****************************************************************************

func TestPUT_Ventures_6(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When an existing Venture is modified and PUT to the server
		Then ensure the response code is 200
		And the 'wrap' query parameter has been specified
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And the body is a JSON object representing a WrappedReply
		And the wrapped data is the updated input Venture
		...`)

	_beginVenTest()
	defer _endVenTest()

	input := v.ModVenture{
		IDs:   "1",
		Props: "description, state, order_ids, extra",
		Values: v.Venture{
			Description: "Black blizzard",
			State:       "In progress",
			OrderIDs:    "1,2,3",
			Extra:       "colour: black; power: 9000",
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures?wrap",
		Method: "PUT",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureHttpMethods)

	_, output := v.AssertWrappedVentureSliceFromReader(t, res.Body)
	require.Len(t, output, 1)

	input.Values.ID = "1"
	input.Values.IsDead = false
	v.AssertGenericVenture(t, output[0])
	v.AssertVentureModEquals(t, input.Values, output[0])
}

// ****************************************************************************
// (OPTIONS) /ventures
// ****************************************************************************

func TestOPTIONS_Ventures(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When /ventures OPTIONS are requested
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And there is NO response body
		...`)

	_beginVenTest()
	defer _endVenTest()

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "OPTIONS",
	}
	res := req.fire()
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertNoContentHeaders(t, res, ventureHttpMethods)
	assertEmptyBody(t, res.Body)
}

// ****************************************************************************
// (?) /ventures
// ****************************************************************************

func TestINVALID_Ventures(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
	 	When /ventures is called using invalid methods
		Then ensure the response code is 405
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, and OPTIONS
		And there is NO response body
		...`)

	_beginVenTest()
	defer _endVenTest()

	verifyNotAllowedMethods(t, "http://localhost:8080/ventures", ventureHttpMethods)
}
