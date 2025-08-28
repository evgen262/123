package news

import (
	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

type categoryPresenter struct{}

func (c *categoryPresenter) NewCategoryToDTO(nc *viewNews.NewCategory) *dtoNews.NewCategory {
	if nc == nil {
		return nil
	}
	return &dtoNews.NewCategory{
		Name:       nc.Name,
		AuthorID:   nc.AuthorID,
		Visibility: c.VisibilityToDTO(nc.Visibility),
	}
}

func (p *categoryPresenter) CategoryToView(c *entityNews.Category) *viewNews.Category {
	if c == nil {
		return nil
	}
	return &viewNews.Category{
		ID:         c.ID,
		Name:       c.Name,
		Visibility: p.VisibilityToView(&c.Visibility),
		UpdatedAt:  c.UpdatedAt,
	}
}

func (p *categoryPresenter) VisibilityToDTO(c *viewNews.CategoryVisibility) *dtoNews.CategoryVisibility {
	if c == nil {
		return nil
	}

	return &dtoNews.CategoryVisibility{
		Condition:        c.Condition,
		ComplexIDs:       c.ComplexIDs,
		OIVs:             c.OIVs,
		OrgIDs:           c.OrgIDs,
		ProductIDs:       c.ProductIDs,
		SubdivisionNames: c.SubdivisionNames,
		PositionNames:    c.PositionNames,
		EmployeeIDs:      c.EmployeeIDs,
		RoleNames:        c.RoleNames,
	}
}

func (p *categoryPresenter) VisibilityToView(c *entityNews.CategoryVisibility) *viewNews.CategoryVisibility {
	if c == nil {
		return nil
	}
	return &viewNews.CategoryVisibility{
		Condition:        c.Condition,
		ComplexIDs:       c.ComplexIDs,
		OIVs:             c.OIVs,
		OrgIDs:           c.OrgIDs,
		ProductIDs:       c.ProductIDs,
		SubdivisionNames: c.SubdivisionNames,
		PositionNames:    c.PositionNames,
		EmployeeIDs:      c.EmployeeIDs,
		RoleNames:        c.RoleNames,
	}
}

func (p *categoryPresenter) UpdateCategoryToDTO(c *viewNews.UpdateCategory) *dtoNews.UpdateCategory {
	if c == nil {
		return nil
	}
	return &dtoNews.UpdateCategory{
		ID:         c.ID,
		Name:       c.Name,
		UpdatedAt:  c.UpdatedAt,
		AuthorID:   c.AuthorID,
		Visibility: p.VisibilityToDTO(c.Visibility),
	}
}

func (p *categoryPresenter) CategoryToResult(c *entityNews.Category) *viewNews.CategoryResult {
	if c == nil {
		return nil
	}

	return &viewNews.CategoryResult{
		ID:        c.ID,
		UpdatedAt: c.UpdatedAt,
	}
}
