package models

import (
	"gorm.io/gorm"
)

type Santri struct {
	gorm.Model
	Nama         string `json:"nama"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Alamat       string `json:"alamat"`
	NamaWali     string `json:"nama_wali"`
	Telp         string `json:"telp"`
	NIS          string `json:"nis" gorm:"uniqueIndex"`
	NISN         string `json:"nisn" gorm:"uniqueIndex"`
	Kelas        string `json:"kelas"`
	Status       string `json:"status" gorm:"default:Aktif"`
}
