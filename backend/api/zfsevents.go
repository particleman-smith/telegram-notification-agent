/*
Responsd to ZFS event HTTP requests
*/

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Test responsds to the endpoint at /zfs-event/test
*/
func Test(writer http.ResponseWriter, request *http.Request) {
	var returnStr = "ZFS Event 'Test' received."
	fmt.Println(returnStr)
	json.NewEncoder(writer).Encode(returnStr)
}
