package user

import (
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

//func (m MockRepo) GetNotes() *model.HttpResponse {
//	args := m.Called()
//	res, _ := args.Get(0).(*model.HttpResponse)
//	return res
//}
//
//func (m MockRepo) GetNoteByID(id int) *model.HttpResponse {
//	args := m.Called(id)
//	res, _ := args.Get(0).(*model.HttpResponse)
//	return res
//}
//
//func (m MockRepo) Creat(n model.Note) *model.HttpResponse {
//	args := m.Called(n)
//	res, _ := args.Get(0).(*model.HttpResponse)
//	return res
//}
