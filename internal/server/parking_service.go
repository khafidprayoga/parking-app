package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/khafidprayoga/parking-app/internal/types"
)

type ParkingServiceImpl struct {
	mu sync.RWMutex

	lotCapacity int            `json:"lot_capacity"`
	store       []*types.Car   `json:"store"`
	revenue     float64        `json:"revenue"`
	tx          map[string]int `json:"-"`
}

func NewParkingService() *ParkingServiceImpl {
	return &ParkingServiceImpl{
		tx: make(map[string]int),
	}
}

func (p *ParkingServiceImpl) Status() (_ []byte, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	countAllTx := 0
	for _, perCarTxHistory := range p.tx {
		countAllTx += perCarTxHistory
	}

	state := types.AppStatus{
		Revenue:            p.revenue,
		LotParkingCapacity: p.lotCapacity,
		TxCount:            countAllTx,
		CarList:            p.store,
	}

	dataBytes, errMarshall := json.Marshal(state)
	if errMarshall != nil {
		err = fmt.Errorf("failed to marshall parking data")
		return
	}

	// stub: simulating deadline for conn lifetime
	//time.Sleep(10 * time.Second)

	return dataBytes, nil
}

func (p *ParkingServiceImpl) OpenParkingArea(parkingCap int) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.lotCapacity > 0 {
		err = fmt.Errorf("failed, already initalize the parking lot capacity")
		return
	}

	p.lotCapacity = parkingCap
	p.store = make([]*types.Car, parkingCap)

	return
}

func (p *ParkingServiceImpl) EnterArea(request types.CarDTO) (areaId int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// validate if  car number not already exist on the parking area
	for _, car := range p.store {
		if car != nil {
			if strings.EqualFold(car.PoliceNumber, request.PoliceNumber) {
				err = fmt.Errorf("failed, already parked the parking lot capacity")
				return
			}
		}
	}

	for index, car := range p.store {
		// allocating nearest parking lot from the door gateway
		if car == nil {
			id := index + 1

			p.store[index] = &types.Car{
				Id:           request.RequestId,
				AreaNumber:   id,
				PoliceNumber: request.GetPoliceNumber(),
				ParkingAt:    time.Now(),
				ExitAt:       nil,
			}

			areaId = id
			return areaId, nil
		}
	}

	// default state when loop is not returned immediately
	err = fmt.Errorf("parking lot capacity is full")
	return
}

func (p *ParkingServiceImpl) LeaveArea(req types.CarDTO) (exitedCar types.Car, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	carDetail := types.Car{}
	carIndex := -1

	for i, car := range p.store {
		if car != nil {
			if strings.EqualFold(car.PoliceNumber, req.PoliceNumber) {
				carDetail = *car
				carIndex = i
				break
			}
		}
	}

	if carIndex < 0 {
		err = fmt.Errorf("car with police number %s does not exist", req.PoliceNumber)
		return
	}

	start := carDetail.ParkingAt
	end := start.Add(time.Duration(req.Hours) * time.Hour)
	carDetail.ExitAt = &end
	carDetail.Cost = p.calculateCost(req.Hours)

	// pay the tx cost
	p.pay(req.GetPoliceNumber())

	// flush
	p.revenue = p.revenue + carDetail.Cost
	p.store[carIndex] = nil

	exitedCar = carDetail
	return
}

func (p *ParkingServiceImpl) pay(policeNumber string) {
	// on existing tx book history
	if val, ok := p.tx[policeNumber]; ok {
		txCount := val + 1
		p.tx[policeNumber] = txCount
		return
	}

	// new member
	p.tx[policeNumber] = 1
}

func (p *ParkingServiceImpl) calculateCost(hours int) float64 {
	const baseCost = 10
	const baseCostMinHours = 2
	if hours <= baseCostMinHours {
		return baseCost
	}

	extraHours := hours - 2

	return float64(baseCost + (extraHours * baseCost))
}
