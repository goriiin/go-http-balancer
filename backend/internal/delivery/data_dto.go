package delivery

import "github.com/goriiin/go-http-balancer/backend/internal/domain"

type dataDto struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

func (d dataDto) toDomain() domain.SomeData {
	return domain.SomeData{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
	}
}

func domainToDto(data domain.SomeData) dataDto {
	return dataDto{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
	}
}
