package pgutils

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// DeviceSession -> device_sessions
type DeviceSession struct {
	ID              int64  //id
	SessionID       string // session_id
	State           string // state
	CreatedTime     string // created_time
	LastUpdatedTime string
}

func (ds DeviceSession) String() string {
	return fmt.Sprintf("DeviceSession<%d %s %s>", ds.ID, ds.SessionID, ds.State)
}

// DeviceResult -> device_results
type DeviceResult struct {
	ID          int64
	SessionID   string
	Model       string
	MacAddress  string
	IpAddress   string
	Netmask     string
	Gateway     string
	Hostname    string
	Kernel      string
	Ap          string
	FirmwareVer string
	Description string
}

func (dr DeviceResult) String() string {
	return fmt.Sprintf("DeviceResult<%d %s %s>", dr.ID, dr.MacAddress, dr.Hostname)
}

// createSchema creates database schema for DeviceResults and DeviceSessions models.
func createDeviceSchema(db *pg.DB) error {
	models := []interface{}{
		(*DeviceSession)(nil),
		(*DeviceResult)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}

	return nil
}
