package db

type DummyDataPage struct {
    Total *int         `json:"total"`
    Data []*Dummy       `json:"data"`
}

func NewEmptyDummyDataPage() *DummyDataPage {
    return &DummyDataPage{}
}
func NewDummyDataPage(total *int, data []*Dummy) *DummyDataPage {
    return &DummyDataPage{ Total: total, Data: data }
}