package models

import (
	"context"
	"donateapp/pkg/entities"
	"log"
	"time"
)

type Donation entities.Donation
type DonationBody entities.DonationBody

func (d *Donation) AddDonations(
	donation DonationBody, userid int,
) (don *Donation, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	q := "INSERT INTO donations (name,photo,location,created_at,updated_at,userid) VALUES (?,?,?,?,?,?)"
	_, err = db.ExecContext(ctx, q, donation.Name, donation.Photo, donation.Location, time.Now(), time.Now(), userid)

	if err != nil {
		return nil, err
	}
	return don, nil
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

func (d *Donation) GetAllDonations() ([]*Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var donations []*Donation

	q := `SELECT * FROM donations`
	rows, err := db.QueryContext(ctx, q)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var donation Donation
		err := rows.Scan(
			&donation.ID,
			&donation.Name,
			&donation.Photo,
			&donation.Location,
			&donation.CreatedAt,
			&donation.UpdatedAt,
			&donation.UserID,
		)

		if err != nil {
			log.Fatalln(err)
		}

		donations = append(donations, &donation)
	}
	return donations, nil
}

func (d *Donation) UpdateDonation(
	id int, userid int, body DonationBody,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	q := `UPDATE donations SET name = ?, photo = ?, location = ?, updated_at = ? WHERE id = ? AND userid = ?`
	_, err := db.ExecContext(
		ctx,
		q,
		body.Name,
		body.Photo,
		body.Location,
		time.Now(),
		id,
		userid,
	)

	if err != nil {
		return "", err
	}

	//_ = row

	return "Updated successfully", nil
}

func (d *Donation) DeleteDonation(id int, userid int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	q := `DELETE FROM donations WHERE id = ? AND userid = ?`
	_, err := db.ExecContext(ctx, q, id, userid)
	//row, err := db.ExecContext(ctx, q, id, userid)

	if err != nil {
		return "", err
	}
	return "Deleted successfully", nil
}
