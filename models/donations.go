package models

import (
	"context"
	"time"
)

type Donation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Photo     string    `json:"photo"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"userid"`
}

func (d *Donation) AddDonations(donation Donation, userid int) (*Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()
	q := "INSERT INTO donations (name,photo,location,created_at,updated_at,userid) VALUES (?,?,?,?,?,?)"

	_, err := db.ExecContext(ctx, q, donation.Name, donation.Photo, donation.Location, time.Now(), time.Now(), userid)

	if err != nil {
		return nil, err
	}

	return &donation, nil
}

func (d *Donation) GetDonationByID(id int) (*Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	q := `SELECT * FROM donations WHERE id = ?`

	row := db.QueryRowContext(ctx, q, id)

	var donation Donation

	err := row.Scan(
		&donation.ID,
		&donation.Name,
		&donation.Photo,
		&donation.Location,
		&donation.CreatedAt,
		&donation.UpdatedAt,
		&donation.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &donation, nil
}
