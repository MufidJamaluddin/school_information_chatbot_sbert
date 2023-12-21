package factory

import "chatbot_be_go/src/persistence/rest/dto"

func CreateTextWAMessage(
	contextMessageId string,
	recipientTo string,
	text string,
	includeUrl bool,
) *dto.WASendMessageDto {
	return &dto.WASendMessageDto{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		RecipientTo:      recipientTo,
		Type:             "text",
		Context: dto.WAContextDto{
			AttachMessageId: contextMessageId,
		},
		Text: dto.WATextMessageDto{
			PreviewUrl: includeUrl,
			Body:       text,
		},
	}
}

func CreateImageWAMessage(
	contextMessageId string,
	recipientTo string,
	imageUrl string,
) *dto.WASendMessageDto {
	return &dto.WASendMessageDto{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		RecipientTo:      recipientTo,
		Type:             "image",
		Context: dto.WAContextDto{
			AttachMessageId: contextMessageId,
		},
		Image: dto.WAImageDto{
			Link: imageUrl,
		},
	}
}

func CreateLocationWAMessage(
	contextMessageId string,
	recipientTo string,
	locationName string,
	locationAddress string,
	latitude float32,
	longitude float32,
) *dto.WASendMessageDto {
	return &dto.WASendMessageDto{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		RecipientTo:      recipientTo,
		Type:             "location",
		Context: dto.WAContextDto{
			AttachMessageId: contextMessageId,
		},
		Location: dto.WALocationDto{
			Longitude: longitude,
			Latitude:  latitude,
			Name:      locationName,
			Address:   locationAddress,
		},
	}
}
