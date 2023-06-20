package blockchain

import "fmt"

// queryResponseError represents the query response error DTO struct
type queryResponseError struct {
	code      string
	message   string
	function  string
	arguments []string
	address   string
}

// NewQueryResponseError creates a new instance of queryResponseError
func NewQueryResponseError(code string, message string, function string, address string, arguments ...string) *queryResponseError {
	return &queryResponseError{
		code:      code,
		message:   message,
		function:  function,
		arguments: arguments,
		address:   address,
	}
}

// Error returns the error string
func (err *queryResponseError) Error() string {
	return fmt.Sprintf("got response code '%s' and message '%s' while querying function '%s' with arguments %v "+
		"and address %s", err.code, err.message, err.function, err.arguments, err.address)
}
