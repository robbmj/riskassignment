
package main

import (
	"fmt"
	"github.com/robbmj/riskassignment/utils"
	"net/http"
)

type myHandler struct {
	response string
}

func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, mh.response)
}

func main() {
	assignment := utils.WriteMyAssignment()
	h := myHandler{ assignment }
    http.ListenAndServe("localhost:4000", h)
}