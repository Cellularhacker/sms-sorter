package data

import "github.com/globalsign/mgo"

//Context is a struct holding the mongo session
type Context struct {
	Session *mgo.Session
}

//Close closes the mongo session
func (c *Context) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}

//SpamDBCollection returns mongodb collection of a given name
func (c *Context) SpamDBCollection(name string) *mgo.Collection {
	return c.Session.DB(SpamDBName).C(name)
}
//SpamDBCollection returns mongodb collection of a given name
func (c *Context) SmsDBCollection(name string) *mgo.Collection {
	return c.Session.DB(SmsDBName).C(name)
}

//NewContext creates a new context and initializes the session
func NewTokyoContext() *Context {
	session := getTokyoSession()
	c := &Context{
		Session: session,
	}
	return c
}
