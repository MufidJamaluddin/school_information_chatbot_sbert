package dto

type WASendMessageDto struct {
	MessagingProduct string           `json:"whatsapp" example:"whatsapp"`
	RecipientType    string           `json:"recipient_type" example:"individual"`
	RecipientTo      string           `json:"to" example:"+6281563532000"`
	Type             string           `json:"type" example:"text"`
	Context          WAContextDto     `json:"context,omitempty"`
	Text             WATextMessageDto `json:"text,omitempty"`
	Image            WAImageDto       `json:"image,omitempty"`
	Location         WALocationDto    `json:"location,omitempty"`
}

type WAContextDto struct {
	AttachMessageId string `json:"message_id,omitempty"`
}

type WATextMessageDto struct {
	PreviewUrl bool   `json:"preview_url,omitempty" example:"false"`
	Body       string `json:"body" example:"Ada mas"`
}

type WAImageDto struct {
	Link string `json:"link" example:"https://awsimages.detik.net.id/community/media/visual/2021/09/11/binus-university-dok-binus.jpeg?w=600"`
}

type WALocationDto struct {
	Longitude float32 `json:"longitude" validate:"gte=-180,lte=180"`
	Latitude  float32 `json:"latitude" validate:"gte=-90,lte=90"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}
