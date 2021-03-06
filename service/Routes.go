package service

import "net/http"

// A Route is a route exposed by the webservice
// The name of the route should be prefaced with the minimal level of authentication required to
// access the endpoint. I.E. Staff_Test can be accessed by Staff or Above. Any route starting
// with _ can be accessed by anyone
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is an array of multiple Route objects
type Routes []Route

var routes = Routes{
	Route{
		"_Test",
		"GET",
		"/",
		HeartBeat,
	},
	Route{
		"Admin_Password_Check",
		"GET",
		"/adminPasswordCheck",
		AdminPasswordCheck,
	},
	Route{
		"Create_Game",
		"POST",
		"/createGame",
		CreateGame,
	},
	Route{
		"Join_Game",
		"GET",
		"/joinGame",
		JoinGame,
	},
	Route{
		"Delete_Game",
		"DELETE",
		"/deleteGame",
		DeleteGame,
	},
	Route{
		"List_Games",
		"GET",
		"/listGames",
		ListGames,
	},
	Route{
		"Game_Info",
		"GET",
		"/gameInfo",
		GameInfo,
	},
	Route{
		"List_Teams",
		"GET",
		"/listTeams",
		ListTeams,
	},
	Route{
		"Team_Info",
		"GET",
		"/teamInfo",
		TeamInfo,
	},
	Route{
		"Start_Game",
		"POST",
		"/startGame",
		StartGame,
	},
}
