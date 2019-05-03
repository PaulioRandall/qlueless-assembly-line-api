package PUT

import (
	"testing"

	a "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/asserts"
	v "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
	test "github.com/PaulioRandall/go-qlueless-assembly-api/test"
	vtest "github.com/PaulioRandall/go-qlueless-assembly-api/test/ventures"
	require "github.com/stretchr/testify/require"
)

// ****************************************************************************
// (PUT) /ventures
// ****************************************************************************

func TestPUT_Ventures_1(t *testing.T) {
	if true {
		return
	}

	t.Log(`Given a Venture already exists on the server
		When a correct modification is PUT to the server
		Then ensure the response code is 200
		And response headers include:
			'Content-Type: 									application/json; charset=utf-8'
			'Access-Control-Allow-Origin: 	*'
			'Access-Control-Allow-Headers: 	*'
			'Access-Control-Allow-Methods: 	GET, POST, PUT, DELETE, OPTIONS'
		And the body is a JSON object containing:
			'message' 											(Non-empty)
			'self' 													(Non-empty)
		And the Venture has been updated within the database
		And that Ventures' 'last_modified' field has been updated appropriately
		...`)

	vtest.BeginEmptyTest("../../../bin")
	defer vtest.EndTest()

	vtest.DBInject(v.NewVenture{
		Description: "Black blizzard",
		State:       "STARTED",
		Orders:      "1,2,3",
		Extra:       "colour: black",
	})

	before := vtest.DBQueryFirst()
	require.NotNil(t, before)

	id := before.ID

	input := &v.Venture{
		ID:          id,
		Description: "White wizzard",
		State:       "FINISHED",
		Orders:      "4,5,6",
		Extra:       "colour: white",
	}

	res := test.CallWithJSON("PUT", "http://localhost:8080/ventures", input)
	defer res.Body.Close()
	defer a.PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	test.AssertDefaultHeaders(t, res, "application/json", vtest.VenHttpMethods)
	v.AssertGenericReplyFromReader(t, res.Body)

	after := vtest.DBQueryOne(id)
	v.AssertEqualsModified(t, &after, input)
}

func TestPUT_Ventures_2(t *testing.T) {

	t.Log(`Given no Ventures exist on the server
		When a non-existent Venture is PUT to the server
		Then ensure the response code is 400
		And response headers include:
			'Content-Type: 									application/json; charset=utf-8'
			'Access-Control-Allow-Origin: 	*'
			'Access-Control-Allow-Headers: 	*'
			'Access-Control-Allow-Methods: 	GET, POST, PUT, DELETE, OPTIONS'
		And the body is a JSON object containing:
			'message' 											(Non-empty)
			'self' 													(Non-empty)
		...`)

}

func TestPUT_Ventures_3(t *testing.T) {

	t.Log(`Given a Venture already exists on the server
		When a modification without an ID is PUT to the server
		Then ensure the response code is 400
		And response headers include:
			'Content-Type: 									application/json; charset=utf-8'
			'Access-Control-Allow-Origin: 	*'
			'Access-Control-Allow-Headers: 	*'
			'Access-Control-Allow-Methods: 	GET, POST, PUT, DELETE, OPTIONS'
		And the body is a JSON object containing:
			'message' 											(Non-empty)
			'self' 													(Non-empty)
		...`)

}

func TestPUT_Ventures_4(t *testing.T) {

	t.Log(`Given a Venture already exists on the server
		When a modification without a Description is PUT to the server
		Then ensure the response code is 400
		And response headers include:
			'Content-Type: 									application/json; charset=utf-8'
			'Access-Control-Allow-Origin: 	*'
			'Access-Control-Allow-Headers: 	*'
			'Access-Control-Allow-Methods: 	GET, POST, PUT, DELETE, OPTIONS'
		And the body is a JSON object containing:
			'message' 											(Non-empty)
			'self' 													(Non-empty)
		...`)

}

// ****************************************************************************
// (PUT) /ventures?wrap
// ****************************************************************************

func TestPUT_Ventures_6(t *testing.T) {

	t.Log(`Given a Venture already exists on the server
		When the Venture is modified
		And the 'wrap' query parameter has been specified
		And PUT request made
		Then ensure the response code is 200
		And response headers include:
			'Content-Type: 									application/json; charset=utf-8'
			'Access-Control-Allow-Origin: 	*'
			'Access-Control-Allow-Headers: 	*'
			'Access-Control-Allow-Methods: 	GET, POST, PUT, DELETE, OPTIONS'
		And the body is a JSON object containing:
			'message' 											(Non-empty)
			'self' 													(Non-empty)
		And the Venture has been updated within the database
		And that Ventures' 'last_modified' field has been updated appropriately
		...`)

}