package main

type Progress struct {
	Percentage float32 `json:"percentage,omitempty"`
	DeviceID   string  `json:"device_id,omitempty"`
	Progress   string  `json:"progress,omitempty"`
	Document   string  `json:"document,omitempty"`
	Device     string  `json:"device,omitempty"`
	Timestamp  int32   `json:"timestamp,omitempty"`
}
