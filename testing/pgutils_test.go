package testing_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Atop-NMS-team/pgutils"
)

func TestDB(t *testing.T) {
	c, err := pgutils.NewClient()

	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close()
	//
	pgdb, err := c.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS device_sessions`)
	if err != nil {
		t.Error(err)
	}
	//
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})

	err = c.Insert(&pgutils.DeviceSession{
		SessionID:       "123",
		State:           "running",
		CreatedTime:     time.Now().String(),
		LastUpdatedTime: time.Now().String(),
	})
	if err != nil {
		t.Error(err)
	}

	c.Insert(&pgutils.DeviceSession{
		SessionID:       "alantroll",
		State:           "running",
		CreatedTime:     time.Now().String(),
		LastUpdatedTime: time.Now().String(),
	})
	if err != nil {
		t.Error(err)
	}

	// var query = &pgutils.DeviceSession{
	// 	State: "running",
	// }
	var results = []pgutils.DeviceSession{}
	err = c.Query(&results, pgutils.QueryExpr{Expr: "state = ?", Value: "running"})

	if err != nil {
		t.Error(err)
	}
	if len(results) != 2 {
		t.Error("slice should has 2 elem but ", len(results))
	}

	update := &pgutils.DeviceSession{
		ID:        1,
		SessionID: "updated-1",
		State:     "done",
	}
	err = c.Update(update)
	if err != nil {
		t.Error(err)
	}

	err = c.Query(&results, pgutils.QueryExpr{Expr: "id = ?", Value: 1})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Error("Results len should be 1 but ", len(results))
	}
	if results[0].ID != 1 {
		t.Error("Result's ID should be 1 but ", results[0].ID)
	}
	if results[0].SessionID != "updated-1" {
		t.Error("Result's SessionID should be updated-1 but ", results[0].SessionID)
	}

}

func TestInsertMany(t *testing.T) {
	// create new client
	c, err := pgutils.NewClient()

	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close()
	pgdb, err := c.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS device_sessions`)
	if err != nil {
		t.Error(err)
	}
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})

	var sessions = []pgutils.DeviceSession{
		{
			SessionID:       "1",
			State:           "running",
			CreatedTime:     time.Now().String(),
			LastUpdatedTime: time.Now().String(),
		},
		{
			SessionID:       "2",
			State:           "running",
			CreatedTime:     time.Now().String(),
			LastUpdatedTime: time.Now().String(),
		},
	}
	err = c.Insert(&sessions)
	if err != nil {
		t.Error("Insert many error ", err)
	}
	var results = []pgutils.DeviceSession{}
	c.Query(&results, pgutils.QueryExpr{Expr: "state = ?", Value: "running"})
	if len(results) != 2 {
		t.Error("Should have two results ")
	}
	t.Log(results)
}

func updateSessionState(ori string, newstate string) string {
	data := map[string]string{}
	states := strings.Split(ori, "|")
	for _, i := range states {
		kv := strings.Split(i, ":")
		data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	data["snmp"] = newstate
	return fmt.Sprintf("gwd:%s|snmp:%s", data["gwd"], data["snmp"])
}

func TestRealCase(t *testing.T) {

	c, err := pgutils.NewClient()

	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close()
	pgdb, err := c.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS device_sessions`)
	if err != nil {
		t.Error(err)
	}
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})

	// Scan service init
	err = c.Insert(&pgutils.DeviceSession{
		SessionID:       "0001",
		State:           "gwd:running|snmp:running",
		CreatedTime:     time.Now().String(),
		LastUpdatedTime: time.Now().String(),
	})

	// gwd update 0001
	q := []pgutils.DeviceSession{}
	err = c.Query(&q, pgutils.QueryExpr{Expr: "session_id = ?", Value: "0001"})
	if err != nil {
		t.Error(err)
	}
	if len(q) != 1 {
		t.Error(fmt.Errorf("should has only one but %d ", len(q)))
	}

	newState := q[0]
	s := updateSessionState(newState.State, "fail")
	newState.State = s
	newState.LastUpdatedTime = time.Now().String()

	err = c.Update(&newState)
	if err != nil {
		t.Error(err)
	}

}

func TestDeviceResult(t *testing.T) {
	c, err := pgutils.NewClient()

	if err != nil {
		t.Error(err)
		return
	}
	defer c.Close()
	pgdb, err := c.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS device_results`)
	if err != nil {
		t.Error(err)
	}
	c.CreateTable(&pgutils.DeviceResult{}, pgutils.CreateTableOpt{IfNotExists: true})

	devices := []pgutils.DeviceResult{
		{
			SessionID:   "s1",
			Model:       "ECH1001",
			MacAddress:  "00-00-00-00-00-01",
			IpAddress:   "10.0.50.1",
			Netmask:     "255.255.255.0",
			Gateway:     "10.0.50.254",
			Hostname:    "switch",
			Kernel:      "K11.01",
			Ap:          "AP11.01",
			FirmwareVer: "F11.01",
			Description: "atop switch",
		},
		{
			SessionID:   "s1",
			Model:       "ECH1001",
			MacAddress:  "00-00-00-00-00-02",
			IpAddress:   "10.0.50.2",
			Netmask:     "255.255.255.0",
			Gateway:     "10.0.50.254",
			Hostname:    "switch",
			Kernel:      "K11.01",
			Ap:          "AP11.01",
			FirmwareVer: "F11.01",
			Description: "atop switch",
		},
		{
			SessionID:   "s1",
			Model:       "ECH1001",
			MacAddress:  "00-00-00-00-00-03",
			IpAddress:   "10.0.50.3",
			Netmask:     "255.255.255.0",
			Gateway:     "10.0.50.254",
			Hostname:    "switch",
			Kernel:      "K11.01",
			Ap:          "AP11.01",
			FirmwareVer: "F11.01",
			Description: "atop switch",
		},
		{
			SessionID:   "s2",
			Model:       "ECH1001",
			MacAddress:  "00-00-00-00-00-01",
			IpAddress:   "10.0.50.1",
			Netmask:     "255.255.255.0",
			Gateway:     "10.0.50.254",
			Hostname:    "switch",
			Kernel:      "K11.01",
			Ap:          "AP11.01",
			FirmwareVer: "F11.01",
			Description: "atop switch",
		},
		{
			SessionID:   "s2",
			Model:       "ECH1001",
			MacAddress:  "00-00-00-00-00-02",
			IpAddress:   "10.0.50.2",
			Netmask:     "255.255.255.0",
			Gateway:     "10.0.50.254",
			Hostname:    "switch",
			Kernel:      "K11.01",
			Ap:          "AP11.01",
			FirmwareVer: "F11.01",
			Description: "atop switch",
		},
	}

	err = c.Insert(&devices)
	if err != nil {
		t.Error(err)
	}

	var q []pgutils.DeviceResult
	err = c.Query(&q,
		pgutils.QueryExpr{Expr: "session_id = ? ", Value: "s1"},
		pgutils.QueryExpr{Expr: "mac_address = ? ", Value: "00-00-00-00-00-02"},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log("q length ", len(q))
	for _, v := range q {
		jsRet, err := json.MarshalIndent(v, "", "")
		if err != nil {
			t.Error(err)
		}
		t.Log(string(jsRet))
	}
}
