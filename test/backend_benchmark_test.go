package test

import (
	"fmt"
	"testing"

	"github.com/khafidprayoga/parking-app/contract"
	"github.com/khafidprayoga/parking-app/internal/backend"
	"github.com/khafidprayoga/parking-app/internal/types"
)

// setupParkingService adalah helper function untuk setup service
func setupParkingService(useCase contract.IParkingUseCase, capacity int) {
	useCase.OpenParkingArea(capacity)
}

func BenchmarkParkingUseCase_EnterArea(b *testing.B) {
	// Setup dengan dependency injection
	useCase := getParkingUseCase() // Implementasi ini akan disediakan oleh aplikasi
	setupParkingService(useCase, 1000)

	// Reset timer sebelum benchmark
	b.ResetTimer()

	// Benchmark loop
	for i := 0; i < b.N; i++ {
		car := types.CarDTO{
			RequestId:    fmt.Sprintf("req-%d", i),
			PoliceNumber: fmt.Sprintf("B%d", i),
		}
		useCase.EnterArea(car)
	}
}

func BenchmarkParkingUseCase_LeaveArea(b *testing.B) {
	// Setup dengan dependency injection
	useCase := getParkingUseCase()
	setupParkingService(useCase, 1000)

	// Pre-fill parking area
	for i := 0; i < 1000; i++ {
		car := types.CarDTO{
			RequestId:    fmt.Sprintf("req-%d", i),
			PoliceNumber: fmt.Sprintf("B%d", i),
		}
		useCase.EnterArea(car)
	}

	// Reset timer sebelum benchmark
	b.ResetTimer()

	// Benchmark loop
	for i := 0; i < b.N; i++ {
		car := types.CarDTO{
			RequestId:    fmt.Sprintf("req-%d", i),
			PoliceNumber: fmt.Sprintf("B%d", i%1000), // Cycle through existing cars
			Hours:        2,
		}
		useCase.LeaveArea(car)
	}
}

func BenchmarkParkingUseCase_EnterAndLeave(b *testing.B) {
	// Setup dengan dependency injection
	useCase := getParkingUseCase()
	setupParkingService(useCase, 1000)

	// Reset timer sebelum benchmark
	b.ResetTimer()

	// Benchmark loop
	for i := 0; i < b.N; i++ {
		// Enter
		enterCar := types.CarDTO{
			RequestId:    fmt.Sprintf("req-%d", i),
			PoliceNumber: fmt.Sprintf("B%d", i),
		}
		useCase.EnterArea(enterCar)

		// Leave
		leaveCar := types.CarDTO{
			RequestId:    fmt.Sprintf("req-%d", i),
			PoliceNumber: fmt.Sprintf("B%d", i),
			Hours:        2,
		}
		useCase.LeaveArea(leaveCar)
	}
}

func BenchmarkParkingUseCase_Parallel(b *testing.B) {
	// Setup dengan dependency injection
	useCase := getParkingUseCase()
	setupParkingService(useCase, 1000)

	// Reset timer sebelum benchmark
	b.ResetTimer()

	// Benchmark loop dengan parallel
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Enter
			enterCar := types.CarDTO{
				RequestId:    fmt.Sprintf("req-%d", i),
				PoliceNumber: fmt.Sprintf("B%d", i),
			}
			useCase.EnterArea(enterCar)

			// Leave
			leaveCar := types.CarDTO{
				RequestId:    fmt.Sprintf("req-%d", i),
				PoliceNumber: fmt.Sprintf("B%d", i),
				Hours:        2,
			}
			useCase.LeaveArea(leaveCar)
			i++
		}
	})
}

// getParkingUseCase adalah helper function untuk mendapatkan implementasi IParkingUseCase
func getParkingUseCase() contract.IParkingUseCase {
	return backend.NewParkingService()
}
