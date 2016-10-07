package handlers

// Configuration holds the values exchanged between the client and the server
// related to the website configuration file.
type Configuration struct {
	FieldNames  []string
	FieldValues []string
}
