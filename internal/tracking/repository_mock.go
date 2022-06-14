package tracking

import (
	"github.com/stretchr/testify/mock"
	myError "notifications/internal/error"
)

type MockRepo struct {
	mock.Mock
}

func (m MockRepo) FindAll(code string) ([]Shipment, *myError.MyError) {
	args := m.Called(code)
	res, _ := args[0].([]Shipment)
	err, _ := args[1].(*myError.MyError)
	return res, err
}
