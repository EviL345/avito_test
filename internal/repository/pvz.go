package repository

import (
	"database/sql"
	"errors"
	"github.com/EviL345/avito_test/internal/model/dto/response"
	"github.com/EviL345/avito_test/internal/model/entity"
	"github.com/google/uuid"
	"log"
	"time"
)

var (
	ErrPvzNotFound = errors.New("pvz not found")
)

type PVZRepository struct {
	db *sql.DB
}

func NewPVZRepository(db *sql.DB) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (r *PVZRepository) CreatePvz(pvz *entity.Pvz) (*entity.Pvz, error) {
	log.SetPrefix("repository.CreatePvz")

	query := `INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`

	pvz.Id = uuid.New()
	pvz.RegistrationDate = time.Now()

	if _, err := r.db.Exec(query, pvz.Id, pvz.RegistrationDate, pvz.City); err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return pvz, nil
}

func (r *PVZRepository) GetPvz(startDate, endDate *time.Time, page, limit *int) ([]response.PvzInfo, error) {
	log.SetPrefix("repository.PvzInfo")

	query := `
        WITH filtered_receptions AS (
            SELECT 
                r.pvz_id
            FROM 
                reception r
            WHERE 
                ($1::timestamp IS NULL OR r.reception_datetime >= $1) AND 
                ($2::timestamp IS NULL OR r.reception_datetime <= $2)
            GROUP BY 
                r.pvz_id
        )
        SELECT 
            p.id, p.city, p.registration_date,
            r.id, r.reception_datetime, r.status, r.pvz_id,
            pr.id, pr.acceptance_datetime, pr.product_type, pr.reception_id
        FROM 
            pvz p
        JOIN 
            filtered_receptions fr ON p.id = fr.pvz_id
        LEFT JOIN 
            reception r ON p.id = r.pvz_id
        LEFT JOIN 
            product pr ON r.id = pr.reception_id
        ORDER BY 
            p.id, r.reception_datetime, pr.acceptance_datetime
        LIMIT $3 OFFSET $4
    `

	offset := (*page - 1) * (*limit)

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		log.Printf("error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	type tempPvzInfo struct {
		pvzInfo       *response.PvzInfo
		receptionsMap map[uuid.UUID]*response.ReceptionsWithProducts
	}

	pvzMap := make(map[uuid.UUID]*tempPvzInfo)
	result := make([]*response.PvzInfo, 0)
	hasRows := false

	for rows.Next() {
		hasRows = true
		var pvzID, receptionID, receptionPVZID, productID, productReceptionID uuid.UUID
		var pvzCity, receptionStatus, productType sql.NullString
		var pvzRegistrationDate, receptionDateTime sql.NullTime
		var productDateTime sql.NullTime

		err = rows.Scan(
			&pvzID, &pvzCity, &pvzRegistrationDate,
			&receptionID, &receptionDateTime, &receptionStatus, &receptionPVZID,
			&productID, &productDateTime, &productType, &productReceptionID,
		)
		if err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, err
		}

		if _, exists := pvzMap[pvzID]; !exists {
			pvz := response.Pvz{
				Id:               pvzID,
				RegistrationDate: pvzRegistrationDate.Time,
				City:             pvzCity.String,
			}
			pvzInfo := &response.PvzInfo{
				Pvz:        pvz,
				Receptions: []response.ReceptionsWithProducts{},
			}
			pvzMap[pvzID] = &tempPvzInfo{
				pvzInfo:       pvzInfo,
				receptionsMap: make(map[uuid.UUID]*response.ReceptionsWithProducts),
			}
			result = append(result, pvzInfo)
		}

		tempPVZ := pvzMap[pvzID]

		if receptionID != uuid.Nil {
			if _, exists := tempPVZ.receptionsMap[receptionID]; !exists {
				reception := response.Reception{
					Id:       receptionID,
					DateTime: receptionDateTime.Time,
					PvzId:    receptionPVZID,
					Status:   receptionStatus.String,
				}
				receptionWithProducts := response.ReceptionsWithProducts{
					Reception: reception,
					Products:  []response.Product{},
				}
				tempPVZ.pvzInfo.Receptions = append(tempPVZ.pvzInfo.Receptions, receptionWithProducts)
				tempPVZ.receptionsMap[receptionID] = &tempPVZ.pvzInfo.Receptions[len(tempPVZ.pvzInfo.Receptions)-1]
			}

			if productID != uuid.Nil {
				product := response.Product{
					Id:          productID,
					DateTime:    productDateTime.Time,
					Type:        productType.String,
					ReceptionId: productReceptionID,
				}
				tempPVZ.receptionsMap[receptionID].Products = append(
					tempPVZ.receptionsMap[receptionID].Products,
					product,
				)
			}
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("error processing rows: %v", err)
		return nil, err
	}

	if !hasRows {
		return []response.PvzInfo{}, nil
	}

	finalResult := make([]response.PvzInfo, len(result))
	for i, pvzInfo := range result {
		finalResult[i] = *pvzInfo
	}

	return finalResult, nil
}
