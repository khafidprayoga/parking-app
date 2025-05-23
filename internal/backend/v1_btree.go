package backend

import (
	"encoding/json"
	"fmt"
	"github.com/google/btree"
	"sync"
	"time"

	"github.com/khafidprayoga/parking-app/internal/types"
)

type ParkingServiceV1BTree struct {
	mu sync.RWMutex

	lotCapacity int

	store   []*types.Car
	hotspot *btree.BTreeG[int]
	history map[string]int

	revenue float64
	tx      map[string]int
}

func NewParkingServiceBTree() *ParkingServiceV1BTree {
	return &ParkingServiceV1BTree{
		tx: make(map[string]int),
	}
}

func (p *ParkingServiceV1BTree) Status() (_ []byte, err error) {
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

func (p *ParkingServiceV1BTree) OpenParkingArea(parkingCap int) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if parkingCap < 1 {
		err = fmt.Errorf("parking cap must be at least 1")
		return
	}

	if p.lotCapacity > 0 {
		err = fmt.Errorf("failed, already initalize the parking lot capacity")
		return
	}

	p.lotCapacity = parkingCap
	p.store = make([]*types.Car, parkingCap)
	p.hotspot = btree.NewOrderedG[int](32)
	p.history = make(map[string]int)

	for i := 0; i < parkingCap; i++ {
		p.hotspot.ReplaceOrInsert(i)
	}
	return
}

func (p *ParkingServiceV1BTree) EnterArea(request types.CarDTO) (areaId int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(request.PoliceNumber) == 0 {
		err = fmt.Errorf("failed, police number is empty")
		return
	}

	if _, exist := p.history[request.PoliceNumber]; exist {
		err = fmt.Errorf("failed, car already parked")
		return
	}

	// validate if  car number not already exist on the parking area
	openArea, found := p.hotspot.Min()
	if !found {
		err = fmt.Errorf("failed, parking lot is full")
		return
	}
	areaId = openArea + 1

	// for compatible with v1 contract
	in := &types.Car{
		Id:           request.RequestId,
		AreaNumber:   areaId,
		PoliceNumber: request.GetPoliceNumber(),
		ParkingAt:    time.Now(),
		ExitAt:       nil,
	}

	p.hotspot.Delete(openArea)
	p.store[openArea] = in
	p.history[request.PoliceNumber] = openArea

	return
}

func (p *ParkingServiceV1BTree) LeaveArea(req types.CarDTO) (exitedCar types.Car, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if req.Hours < 1 {
		err = fmt.Errorf("failed, parking must be at least 1 hour")
		return
	}

	parkingSpot, exists := p.history[req.PoliceNumber]
	if !exists {
		err = fmt.Errorf("this car %s not exist on parking area", req.PoliceNumber)
		return
	}

	// get the car data
	car := p.store[parkingSpot]

	// free the history mem
	delete(p.history, req.PoliceNumber)
	p.store[parkingSpot] = nil
	p.hotspot.ReplaceOrInsert(parkingSpot)

	start := car.ParkingAt
	end := start.Add(time.Duration(req.Hours) * time.Hour)
	car.ExitAt = &end
	car.Cost = p.calculateCost(req.Hours)

	// pay the tx cost
	p.pay(req.GetPoliceNumber())

	// add revenue
	p.revenue = p.revenue + car.Cost

	exitedCar = *car
	return
}

func (p *ParkingServiceV1BTree) pay(policeNumber string) {
	// on existing tx book history
	if val, ok := p.tx[policeNumber]; ok {
		txCount := val + 1
		p.tx[policeNumber] = txCount
		return
	}

	// new member
	p.tx[policeNumber] = 1
}

func (p *ParkingServiceV1BTree) calculateCost(hours int) float64 {
	const baseCost = 10
	const baseCostMinHours = 2
	if hours <= baseCostMinHours {
		return baseCost
	}

	extraHours := hours - 2

	return float64(baseCost + (extraHours * baseCost))
}
