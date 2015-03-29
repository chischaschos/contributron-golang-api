package main

import "net/http"

func init() {
	//http.HandleFunc("/", doit)

	http.HandleFunc("/pull-people", Wrap(GetPublicMembersList, &MyContext{}))
	http.HandleFunc("/pull-historic-archive", Wrap(GetHistoricArchive, &MyContext{}))
}
