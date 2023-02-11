package sample

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"pcbook/pb"
)

func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	numberOfCores := randomInt(2, 8)
	numberOfThreads := randomInt(numberOfCores, 12)
	minGHz := randomFloat64(2.0, 3.5)
	maxGHz := randomFloat64(minGHz, 5.0)

	return &pb.CPU{
		Brand:         brand,
		Name:          randomCPUName(brand),
		NumberCores:   uint32(numberOfCores),
		NumberThreads: uint32(numberOfThreads),
		MinGhz:        minGHz,
		MaxGhz:        maxGHz,
	}
}

func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)
	minGHz := randomFloat64(1.0, 1.5)
	maxGHz := randomFloat64(minGHz, 2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}

	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGHz,
		MaxGhz: maxGHz,
		Memory: memory,
	}
}

func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_MEGABYTE,
	}
}

func NewSSD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

func NewScreen() *pb.Screen {
	return &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	return &pb.Laptop{
		Id:       randomId(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Memory:   NewRAM(),
		Gpus:     []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   timestamppb.Now(),
	}
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
