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
	}

	db, _ := c.GetDB()

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
	var result = []pgutils.DeviceSession{}
	err = c.Query(&result, "state = ?", "running")
	if err != nil {
		t.Error(err)
	}
	t.Log("q ", result)
	if len(result) != 2 {
		t.Error("slice should has 2 elem ")
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

	err = c.Query(&result, "id = ?", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log("updated result ", result)

	var sessions []pgutils.DeviceSession
	err = db.Model(&sessions).Select()
	if err != nil {
		t.Error(err)
	}
	t.Log("sessions ", sessions)
}
