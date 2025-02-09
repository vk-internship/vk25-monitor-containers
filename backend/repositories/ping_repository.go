package repositories

import (
	"database/sql"
	"fmt"
	"monitoring-backend/models"
)

type PingRepository struct {
	DB *sql.DB
}

func NewPingRepository(db *sql.DB) *PingRepository {
	return &PingRepository{DB: db}
}

func (r *PingRepository) GetAll() ([]models.Ping, error) {
	rows, err := r.DB.Query("SELECT id, ip_address, ping_time, is_success FROM pings")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pings []models.Ping

	for rows.Next() {
		var ping models.Ping

		err := rows.Scan(&ping.ID, &ping.IPAddress, &ping.PingTime, &ping.IsSuccess)

		if err != nil {
			return nil, err
		}

		pings = append(pings, ping)
	}

	return pings, nil
}

func (r *PingRepository) CreateOrUpdate(ping models.Ping) error {
	query := `
        INSERT INTO pings (ip_address, ping_time, is_success)
        VALUES ($1, $2, $3)
        ON CONFLICT (ip_address) DO UPDATE
        SET ping_time = EXCLUDED.ping_time, is_success = EXCLUDED.is_success;
    `

	_, err := r.DB.Exec(query, ping.IPAddress, ping.PingTime, ping.IsSuccess)

	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return nil
}
