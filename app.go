package main

import "net/http"

func init() {
	http.HandleFunc("/pull-people", Wrap(GetPublicMembersList, &MyContext{}))
	http.HandleFunc("/pull-historic-archive", Wrap(GetHistoricArchive, &MyContext{}))
	http.HandleFunc("/pull-current-year-archive", Wrap(GetCurrentYearArchive, &MyContext{}))
	http.HandleFunc("/all-time-stats", Wrap(GetAllTimeStats, &MyContext{}))
}
