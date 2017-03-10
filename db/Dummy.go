package db

type Dummy struct {
    ID      string  `json:"id" bson:"_id"`
    Key     string  `json:"key" bson:"key"`
    Content string  `json:"content" bson:"content"`
}

func NewEmptyDummy() *Dummy {
    return &Dummy {}
}

func NewDummy(id, key, content string) *Dummy {
    return &Dummy { ID: id, Key: key, Content: content }
}

func (c *Dummy) GetID() string {
    return c.ID
}

func (c *Dummy) SetID(id string) {
    c.ID = id
}