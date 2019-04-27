package models

// Configuration is the overall configuration
type Configuration struct {
	Redis RedisConfiguration `json:"redis"`
}

// RedisConfiguration for connection
type RedisConfiguration struct {
	Port     int    `json:"port"`
	Address  string `json:"address"`
	Database int    `json:"db"`
	Password string `json:"password"`
}
