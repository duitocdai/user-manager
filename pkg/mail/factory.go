package mail

import (
	"fmt"

	"github.com/duitocdai/internal/utils"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type ProviderFactory interface {
	CreateProvider(config *MailConfig) (EmailProviderService, error)
}
type MailtrapProviderFactory struct {
}

func (f *MailtrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error) {
	return NewMailTrapProvider(config)
}
func NewProviderFactory(providerType ProviderType) (ProviderFactory, error) {
	switch providerType {
	case ProviderMailtrap:
		return &MailtrapProviderFactory{}, nil
	default:
		return nil, utils.NewError(fmt.Sprintf("unsupported provider type %s", utils.ErrorCode(providerType)), utils.ErrCodeInternal)
	}
}
