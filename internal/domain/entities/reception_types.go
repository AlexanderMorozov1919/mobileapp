package entities

// TherapistData - данные для терапевта
type TherapistData struct {
	BloodPressure string  `json:"blood_pressure" validate:"required"`
	Temperature   float32 `json:"temperature" validate:"required"`
	Anamnesis     string  `json:"anamnesis"`
}

// CardiologistData - данные для кардиолога
type CardiologistData struct {
	ECG        string `json:"ecg"`
	HeartRate  int    `json:"heart_rate" validate:"min=30,max=200"`
	Arrhythmia bool   `json:"arrhythmia"`
}

// NeurologistData - данные для невролога
type NeurologistData struct {
	Reflexes    map[string]string `json:"reflexes"`
	Sensitivity string            `json:"sensitivity"`
	Complaints  []string          `json:"complaints"`
}

// TraumatologistData - данные для травматолога
type TraumatologistData struct {
	InjuryType string `json:"injury_type"`
	XRayResult string `json:"x_ray_result"`
	Fracture   bool   `json:"fracture"`
}
