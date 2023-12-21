package dto

type WebhookInDto struct {
	Object string              `json:"object" example:"whatsapp_business_account"`
	Entry  []WebhookInEntryDto `json:"entry"`
}

type WebhookInEntryDto struct {
	Changes []WebhookInEntryChangeDto `json:"changes"`
}

type WebhookInEntryChangeDto struct {
	Value WebhookInEntryValueDto `json:"value"`
}

type WebhookInEntryValueDto struct {
	Metadata WebhookInEntryMetadataDto  `json:"metadata"`
	Contacts []WebhookInEntryContactDto `json:"contacts"`
	Messages []WebhookInEntryMessageDto `json:"messages"`
}

type WebhookInEntryMetadataDto struct {
	PhoneNoId string `json:"phone_number_id" example:"103055476187368"`
}

type WebhookInEntryContactDto struct {
	Profile WebhookInEntryProfileDto `json:"profile"`
	WaId    string                   `json:"wa_id" example:"6281563532000"`
}

type WebhookInEntryProfileDto struct {
	Name string `json:"name" example:"Mufid Jamaluddin"`
}

type WebhookInEntryMessageDto struct {
	From      string                       `json:"from" example:"6281563532000"`
	Id        string                       `json:"id" example:"103055476187368"`
	Timestamp int64                        `json:"timestamp" example:"1688575667161"`
	Type      string                       `json:"type" example:"text"`
	Text      WebhookInEntryMessageTextDto `json:"text"`
}

type WebhookInEntryMessageTextDto struct {
	Body string `json:"body" example:"Selamat malam, apakah ada yang bisa saya bantu?"`
}
