package pgutils

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type Client struct {
	db *pg.DB
}

func (c *Client) GetDB() (*pg.DB, error) {
	if c.db != nil {
		return c.db, nil
	}
	db, err := connectDB("user", "pass")
	c.db = db
	return c.db, err

}

func (c *Client) Insert(data interface{}) error {
	result, err := c.db.Model(data).Insert()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (c *Client) Query(results interface{}, field string, value interface{}) error {

	err := c.db.Model(results).Where(field, value).Select()
	if err != nil {
		return err
	}
	fmt.Println(results)
	return nil
}

func (c *Client) Update(data interface{}) error {

	_, err := c.db.Model(data).WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}

type DBClient interface {
	GetDB() (*pg.DB, error)

	Query(results interface{}, field string, value interface{}) error
	Insert(model interface{}) error
	Update(update interface{}) error
	// Close()
}

// connectDB
func connectDB(user, pass string) (*pg.DB, error) {
	url := getDBHost(user, pass)
	// opt, err := pg.ParseURL(url)
	// if err != nil {
	// 	log.Errorf("connect to database (%s) was failed err:%v", url, err)
	// 	return nil, err
	// }
	// fmt.Println(opt)
	// disable ssl mode
	// TODO - should enable ssl
	// opt.TLSConfig = nil
	// pgdb := pg.Connect(opt)
	pgdb := pg.Connect(&pg.Options{
		Addr:     url,
		User:     user,
		Password: pass,
	})
	initDB(pgdb)
	// check databas is up
	if err := pgdb.Ping(ctx); err != nil {
		log.Errorf("database not up err:%v", err)
		return nil, err
	}
	return pgdb, nil
}

func initDB(db *pg.DB) {
	if db == nil {
		log.Errorf("got nil db ")
		return
	}
	err := createDeviceSchema(db)
	if err != nil {
		log.Errorf("create schema fail err:%v", err)
		// return
	}
}

// NewClient new PostgreSQL client
func NewClient() (DBClient, error) {
	return NewClientWithAccount("user", "pass")
}

// NewClientWithAccount new PostgreSQL client with account authentication
func NewClientWithAccount(user, pass string) (DBClient, error) {
	pgdb, err := connectDB(user, pass)
	if err != nil {
		return nil, err
	}

	initDB(pgdb)
	return &Client{
		db: pgdb,
	}, nil
}
