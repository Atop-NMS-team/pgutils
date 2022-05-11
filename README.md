# pgutils
PostgreSQL common API 

A common interface to use PostgreSQL DB

## How to use
### Insert
- Insert one, pass structure instance's pointer to `Insert()` funciton
- Insert many, pass slice pointer to `Insert()` funciton
- The parameter's type must create the shema before insert
```go
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
// Insert accept parameter only if its a pointer type
err = c.Insert(&sessions)
if err != nil {
  t.Error("Insert many error ", err)
}
```

### Update