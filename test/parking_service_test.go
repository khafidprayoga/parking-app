package test

import (
	"encoding/json"
	"github.com/khafidprayoga/parking-app/internal/backend"
	"testing"

	"github.com/khafidprayoga/parking-app/contract"
	"github.com/khafidprayoga/parking-app/internal/types"
	"github.com/stretchr/testify/assert"
)

type TestParkingService struct {
	service contract.IParkingUseCase
}

func setupTestParkingService() *TestParkingService {
	return &TestParkingService{
		service: backend.NewParkingServiceBTree(),
	}
}

func TestParkingService_OpenArea(t *testing.T) {
	// Setup
	testService := setupTestParkingService()
	expectedAreaCapacity := 6

	// Test
	err := testService.service.OpenParkingArea(expectedAreaCapacity)

	// Assertions
	assert.NoError(t, err)

	// overriding the current active area will thrown error
	err = testService.service.OpenParkingArea(5)
	assert.Error(t, err)
}

func TestParkingService_EnterArea(t *testing.T) {
	// Setup
	testService := setupTestParkingService()
	expectedAreaCapacity := 6
	err := testService.service.OpenParkingArea(expectedAreaCapacity)
	assert.NoError(t, err)

	// Test cases
	testCases := []struct {
		name         string
		car          types.CarDTO
		expectedArea int
		expectError  bool
	}{
		{
			name: "Successfully park first car",
			car: types.CarDTO{
				RequestId:    "req-1",
				PoliceNumber: "B1234ABC",
			},
			expectedArea: 1,
			expectError:  false,
		},
		{
			name: "Successfully park second car",
			car: types.CarDTO{
				RequestId:    "req-2",
				PoliceNumber: "B5678DEF",
			},
			expectedArea: 2,
			expectError:  false,
		},
		{
			name: "Fail to park duplicate car",
			car: types.CarDTO{
				RequestId:    "req-3",
				PoliceNumber: "B1234ABC", // Duplicate police number
			},
			expectedArea: 0,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			areaId, err := testService.service.EnterArea(tc.car)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArea, areaId)
			}
		})
	}
}

func TestParkingService_LeaveArea(t *testing.T) {
	// Setup
	testService := setupTestParkingService()
	expectedAreaCapacity := 6
	err := testService.service.OpenParkingArea(expectedAreaCapacity)
	assert.NoError(t, err)

	// Park a car first
	car := types.CarDTO{
		RequestId:    "req-1",
		PoliceNumber: "B1234ABC",
	}
	areaId, err := testService.service.EnterArea(car)
	assert.NoError(t, err)
	assert.Equal(t, 1, areaId)

	// Test cases
	testCases := []struct {
		name        string
		car         types.CarDTO
		hours       int
		expectError bool
	}{
		{
			name: "Successfully leave car",
			car: types.CarDTO{
				PoliceNumber: "B1234ABC",
				Hours:        2,
			},
			expectError: false,
		},
		{
			name: "Fail to leave non-existent car",
			car: types.CarDTO{
				PoliceNumber: "B9999XYZ",
				Hours:        2,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exitedCar, err := testService.service.LeaveArea(tc.car)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, exitedCar.ExitAt)
				assert.Equal(t, tc.car.PoliceNumber, exitedCar.PoliceNumber)
				assert.Equal(t, areaId, exitedCar.AreaNumber)
			}
		})
	}
}

func TestParkingService_Status(t *testing.T) {
	// Setup
	testService := setupTestParkingService()
	expectedAreaCapacity := 6
	err := testService.service.OpenParkingArea(expectedAreaCapacity)
	assert.NoError(t, err)

	// Park some cars
	cars := []types.CarDTO{
		{
			RequestId:    "req-1",
			PoliceNumber: "B1234ABC",
		},
		{
			RequestId:    "req-2",
			PoliceNumber: "B5678DEF",
		},
	}

	for _, car := range cars {
		_, err := testService.service.EnterArea(car)
		assert.NoError(t, err)
	}

	// Test
	status, err := testService.service.Status()
	assert.NoError(t, err)
	assert.NotEmpty(t, status)

	// Verify status contains expected information
	var statusData types.AppStatus
	carCount := 0
	err = json.Unmarshal(status, &statusData)

	for _, c := range statusData.CarList {
		if c != nil {
			carCount++
		}
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedAreaCapacity, statusData.LotParkingCapacity)
	assert.Equal(t, 2, carCount)
}

func TestParkingService_EntryPoint(t *testing.T) {
	// Setup
	testService := setupTestParkingService()

	// Test cases untuk OpenParkingArea
	t.Run("OpenParkingArea Edge Cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			capacity    int
			expectError bool
		}{
			{
				name:        "Success open with valid capacity",
				capacity:    6,
				expectError: false,
			},
			{
				name:        "Fail open with zero capacity",
				capacity:    0,
				expectError: true,
			},
			{
				name:        "Fail open with negative capacity",
				capacity:    -1,
				expectError: true,
			},
			{
				name:        "Fail open with already initialized area",
				capacity:    5,
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := testService.service.OpenParkingArea(tc.capacity)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	// Test cases untuk EnterArea
	t.Run("EnterArea Edge Cases", func(t *testing.T) {
		testCases := []struct {
			name         string
			car          types.CarDTO
			expectedArea int
			expectError  bool
		}{
			{
				name: "Success park with valid data",
				car: types.CarDTO{
					RequestId:    "req-1",
					PoliceNumber: "B1234ABC",
				},
				expectedArea: 1,
				expectError:  false,
			},
			{
				name: "Fail park with empty police number",
				car: types.CarDTO{
					RequestId:    "req-2",
					PoliceNumber: "",
				},
				expectedArea: 0,
				expectError:  true,
			},
			{
				name: "Fail park with duplicate police number",
				car: types.CarDTO{
					RequestId:    "req-3",
					PoliceNumber: "B1234ABC",
				},
				expectedArea: 0,
				expectError:  true,
			},
			{
				name: "Success park multiple cars",
				car: types.CarDTO{
					RequestId:    "req-4",
					PoliceNumber: "B5678DEF",
				},
				expectedArea: 2,
				expectError:  false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				areaId, err := testService.service.EnterArea(tc.car)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedArea, areaId)
				}
			})
		}
	})

	// Test cases untuk LeaveArea
	t.Run("LeaveArea Edge Cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			car         types.CarDTO
			expectError bool
		}{
			{
				name: "Success leave with valid data",
				car: types.CarDTO{
					PoliceNumber: "B1234ABC",
					Hours:        2,
				},
				expectError: false,
			},
			{
				name: "Fail leave with non-existent car",
				car: types.CarDTO{
					PoliceNumber: "B9999XYZ",
					Hours:        2,
				},
				expectError: true,
			},
			{
				name: "Fail leave with negative hours",
				car: types.CarDTO{
					PoliceNumber: "B5678DEF",
					Hours:        -1,
				},
				expectError: true,
			},
			{
				name: "Success leave with zero hours",
				car: types.CarDTO{
					PoliceNumber: "B5678DEF",
					Hours:        0,
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				exitedCar, err := testService.service.LeaveArea(tc.car)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if !tc.expectError {
						assert.NotNil(t, exitedCar.ExitAt)
						assert.Equal(t, tc.car.PoliceNumber, exitedCar.PoliceNumber)
					}
				}
			})
		}
	})

	// Test cases untuk Status
	t.Run("Status Edge Cases", func(t *testing.T) {
		// Test status setelah beberapa operasi
		status, err := testService.service.Status()
		assert.NoError(t, err)
		assert.NotEmpty(t, status)

		var statusData types.AppStatus
		err = json.Unmarshal(status, &statusData)
		assert.NoError(t, err)
		assert.Equal(t, 6, statusData.LotParkingCapacity)
		assert.GreaterOrEqual(t, statusData.Revenue, 0.0)
		assert.GreaterOrEqual(t, statusData.TxCount, 0)
	})

	// todo clear the parking lot

	// Test cases untuk Full Capacity
	t.Run("Full Capacity Edge Cases", func(t *testing.T) {
		// Coba parkir mobil sampai penuh
		for i := 3; i <= 7; i++ {
			car := types.CarDTO{
				RequestId:    "req-" + string(rune(i)),
				PoliceNumber: "B" + string(rune(i)) + "XYZ",
			}
			_, err := testService.service.EnterArea(car)
			assert.NoError(t, err)
		}

		// Coba parkir mobil saat penuh
		car := types.CarDTO{
			RequestId:    "req-full",
			PoliceNumber: "B9999XYZ",
		}
		areaId, err := testService.service.EnterArea(car)
		assert.Error(t, err)
		assert.Equal(t, 0, areaId)
	})

	// Test cases untuk Concurrent Operations
	t.Run("Concurrent Operations", func(t *testing.T) {
		// Reset service
		testService = setupTestParkingService()
		err := testService.service.OpenParkingArea(6)
		assert.NoError(t, err)

		// Test concurrent enter operations
		done := make(chan bool)
		for i := 0; i < 3; i++ {
			go func(id int) {
				car := types.CarDTO{
					RequestId:    "req-concurrent-" + string(rune(id)),
					PoliceNumber: "B" + string(rune(id)) + "CONC",
				}
				_, err := testService.service.EnterArea(car)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Verify final state
		status, err := testService.service.Status()
		assert.NoError(t, err)
		carCount := 0
		var statusData types.AppStatus
		err = json.Unmarshal(status, &statusData)
		assert.NoError(t, err)

		for _, c := range statusData.CarList {
			if c != nil {
				carCount++
			}
		}

		assert.Equal(t, 3, carCount)
	})
}
