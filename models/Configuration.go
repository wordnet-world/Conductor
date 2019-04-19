package models

// Configuration is the overall configuration
type Configuration struct {
	CouchDB     CouchDBConfiguration     `json:"couchDB"`
	RoleService RoleServiceConfiguration `json:"roleService"`
}

// CouchDBConfiguration represents a configuration for
// couchDB
type CouchDBConfiguration struct {
	DatabasePort     int    `json:"databasePort"`
	ConnectionString string `json:"connectionString"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Database         string `json:"database"`
}

// RoleServiceConfiguration represents the configuration
// for the role service
type RoleServiceConfiguration struct {
	URL string `json:"url"`
}
