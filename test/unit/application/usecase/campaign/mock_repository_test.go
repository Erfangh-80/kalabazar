package campaign_test

import (
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

type mockRepository struct {
	campaigns map[int64]*domaincampaign.Campaign
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		campaigns: make(map[int64]*domaincampaign.Campaign),
	}
}

func (m *mockRepository) Save(c *domaincampaign.Campaign) error {
	m.campaigns[c.ID] = c
	return nil
}

func (m *mockRepository) FindByID(id int64) (*domaincampaign.Campaign, error) {
	c, ok := m.campaigns[id]
	if !ok {
		return nil, domaincampaign.ErrCampaignNotFound
	}
	return c, nil
}

func (m *mockRepository) Update(c *domaincampaign.Campaign) error {
	m.campaigns[c.ID] = c
	return nil
}
