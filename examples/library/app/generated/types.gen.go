// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import ()

// BooksListRequest
type BooksListRequestMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Correlation ID set by client
		CorrelationId string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		// Genre
		Genre string `json:"genre"`
	}
}

// BooksListResponse
type BooksListResponseMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Correlation ID set by client on corresponding request
		CorrelationId string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		// Books list
		Books []Book `json:"books"`
	}
}

// Book Information
type Book struct {
	// Title
	Title string `json:"title"`
}
