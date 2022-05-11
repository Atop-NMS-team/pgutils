package testing_test

import (
	"pgutils"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
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
	err = c.Query(&results, "state = ?", "running")
	defer c.Close()
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

	err = c.Query(&results, "id = ?", 1)
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
	c.Query(&results, "state = ?", "running")
	if len(results) != 2 {
		t.Error("Should have two results ")
	}
	t.Log(results)
}
