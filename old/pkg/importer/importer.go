package importer

import (
	"github.com/ashikhman/scraper/old/pkg/db"
	"github.com/rs/zerolog/log"
)

type Importer struct {
	domain *db.Domain
	urls   []string
}

func New(domainName string) *Importer {
	domain := db.FindDomainByName(domainName)
	if nil == domain {
		domain = &db.Domain{
			Name: domainName,
		}
		db.SaveDomain(domain)
	}

	return &Importer{
		domain: domain,
		urls:   []string{},
	}
}

func (i *Importer) AddUrl(url string) {
	i.urls = append(i.urls, url)
}

func (i *Importer) Flush() error {
	log.Info().Msg("Flushing Importer")

	err := db.InsertQueueRecords(i.domain.ID, i.urls)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to flush importer")
	}

	log.Info().Msg("Flushing done")

	return nil
}
