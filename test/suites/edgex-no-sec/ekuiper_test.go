package test

import (
	"edgex-snap-testing/test/utils"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type Reading struct {
	TotalCount int `json:"totalCount"`
}

func TestRulesEngine(t *testing.T) {
	t.Cleanup(func() {
		utils.SnapStop(t,
			ekuiperService,
			deviceVirtualSnap,
			ascSnap)
	})

	utils.SnapStart(t,
		ekuiperService,
		deviceVirtualSnap,
		ascSnap)

	t.Run("create stream and rule", func(t *testing.T) {
		utils.Exec(t, `edgex-ekuiper.kuiper-cli create stream stream1 '()WITH(FORMAT="JSON",TYPE="edgex")'`)

		utils.Exec(t,
			`edgex-ekuiper.kuiper-cli create rule rule_edgex_message_bus '
			{
			   "sql":"SELECT * from stream1",
			   "actions": [
				  {
					 "edgex": {
						"connectionSelector": "edgex.redisMsgBus",
						"topicPrefix": "edgex/events/device", 
						"messageType": "request",
						"deviceName": "device-test"
					 }
				  }
			   ]
			}'`)

		// wait device-virtual to come online and produce readings
		utils.WaitServiceOnline(t, 60, deviceVirtualPort)
		utils.TestDeviceVirtualReading(t)

		var reading Reading
		resp, err := http.Get("http://localhost:59880/api/v2/reading/device/name/device-test")
		if err != nil {
			fmt.Print(err)
			return
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&reading); err != nil {
			t.Fatal(err)
		}

		require.Greaterf(t, reading.TotalCount, 0, "No readings have been re-published to EdgeX message bus by ekuiper")
	})
}