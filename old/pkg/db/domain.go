package db

import (
	"context"
	"github.com/rs/zerolog/log"
)

type Domain struct {
	ID   int
	Name string
}

func SaveDomain(domain *Domain) {
	var err error
	if 0 == domain.ID {
		err = pool.QueryRow(context.Background(), "INSERT INTO domain (name) VALUES ($1) RETURNING id", domain.Name).Scan(&domain.ID)
	} else {
		_, err = pool.Exec(context.Background(), "UPDATE domain SET name = $2 WHERE id = $1", domain.ID, domain.Name)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to insert new domain")
	}
}

func FindDomainByName(name string) *Domain {
	var domain Domain

	rows, err := pool.Query(context.Background(), "SELECT id, name FROM domain WHERE name = $1", name)
	if err != nil {
		log.Fatal().Err(err).Str("domain", name).Msg("Failed to fetch the domain")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&domain.ID, &domain.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("Scan() to domain has failed.")
		}

		return &domain
	}

	return nil
}
