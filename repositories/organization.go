package repositories

import (
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	OrganizationRepo = &OrganizationRepository{db}
	return OrganizationRepo
}

func (or *OrganizationRepository) CreateOrganization(organization *models.Organization) (*models.Organization, error) {
	if err := or.db.Create(organization).Error; err != nil {
		return nil, err
	}
	return organization, nil
}

func (or *OrganizationRepository) GetOrganizationByName(name string) *models.Organization {
	var organization models.Organization
	or.db.Where("name = ?", name).First(&organization)
	if organization.ID == 0 {
		return nil
	}
	return &organization
}

func (or *OrganizationRepository) GetOrganizationByID(id uint) (*models.Organization, error) {
	var organization models.Organization
	if err := or.db.Where("id = ?", id).First(&organization).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (or *OrganizationRepository) UpdateOrganization(organization *models.Organization) error {
	return or.db.Save(organization).Error
}

func (or *OrganizationRepository) DeleteOrganization(organization *models.Organization) error {
	return or.db.Delete(organization).Error
}

func (or *OrganizationRepository) GetOrganizations() ([]models.Organization, error) {
	var organizations []models.Organization
	if err := or.db.Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}
