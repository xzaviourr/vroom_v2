package main

import "fmt"

// Stores properties of a function variant
type Variant struct {
	Id             string  // Unique id for the variant
	TaskId         string  `json:"task-identifier"` // Unique task id
	GpuMemory      int64   `json:"gpu-memory"`      // GPU memory (in MB)
	GpuCores       int64   `json:"gpu-cores"`       // Percentage of GPU cores
	Image          string  `json:"image"`           // Docker image file
	StartupLatency float32 `json:"startup-latency"` // Time to start the container
	MinLatency     float32 `json:"min-latency"`     // Latency for request when load is minimum
	MeanLatency    float32 `json:"mean-latency"`    // Latency for request of when load is at average
	MaxLatency     float32 `json:"max-latency"`     // Latency for request when load is at its peak
	Accuracy       float32 `json:"accuracy"`        // Model accuracy
	BatchSize      int64   `json:"batch-size"`      // Batch size
	EndPoint       string  `json:"end-point"`       // Url endpoint extension to access the service
	Port           int64   `json:"port"`            // Port at which service is running internally
	Capacity       float32 `json:"capacity"`        // Maximum load capacity
}

func (v *Variant) String() string {
	return fmt.Sprintf(
		"Function Variant\nId: %s\nTask Id: %s\nGPU Memory: %d\nGPU Cores: %d\nImage: %s\n"+
			"Startup Latency: %f\nMin Latency: %f\nMean Latency: %f\nMax Latency: %f\nAccuracy: %f\n"+
			"Batch Size: %d\nEnd Point: %s\nPort: %d\nCapacity: %f\n",
		v.Id, v.TaskId, v.GpuMemory, v.GpuCores, v.Image, v.StartupLatency, v.MinLatency,
		v.MeanLatency, v.MaxLatency, v.Accuracy, v.BatchSize, v.EndPoint, v.Port, v.Capacity,
	)
}

type VariantStore struct {
	Variants        map[string]*Variant // VariantId -> Variant Info structure
	DatabaseManager *DatabaseManager    // Database handler
}

func initVariantStore() *VariantStore {
	databaseManager := initDatabaseManager()

	variantStore := VariantStore{
		Variants:        databaseManager.loadAllVariantsFromDb(),
		DatabaseManager: databaseManager,
	}
	return &variantStore
}

func (vs *VariantStore) String() string {
	result := "Variant Store\n"
	for _, variant := range vs.Variants {
		result += variant.String() + "\n\n"
	}
	return result
}

func (vs *VariantStore) addVariant(variant *Variant) {
	vs.Variants[variant.Id] = variant
	vs.DatabaseManager.insertVariantInDb(variant)
}

func (vs *VariantStore) removeVariant(variantId string) {
	delete(vs.Variants, variantId)
	// vs.DatabaseManager.removeVariantInDb(variantId)
}

func (vs *VariantStore) getVariant(variantId string) *Variant {
	return vs.Variants[variantId]
}

func (vs *VariantStore) getRelevantVariants(taskId string, accuracy float32,
	latency float32) []*Variant {
	var relevantVariants []*Variant

	// Find the variants that satisfies the constraints
	for _, variant := range vs.Variants {
		if variant.TaskId != taskId {
			continue
		}

		if variant.Accuracy >= accuracy && variant.MinLatency <= latency &&
			variant.MaxLatency >= latency {
			relevantVariants = append(relevantVariants, variant)
		}
	}
	return relevantVariants
}

func (vs *VariantStore) getTaskVariants(taskId string) []*Variant {
	var relevantVariants []*Variant

	// Find the variants that satisfies the constraints
	for _, variant := range vs.Variants {
		if variant.TaskId != taskId {
			continue
		}
		relevantVariants = append(relevantVariants, variant)
	}
	return relevantVariants
}
