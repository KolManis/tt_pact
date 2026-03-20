package converter

import (
	"github.com/KolManis/tt_pact/internal/model"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
)

func SessionToProtoResponse(session *model.Session) *telegramV1.CreateSessionResponse {
	return &telegramV1.CreateSessionResponse{
		SessionId: session.ID,
		QrCode:    session.QRCode,
	}
}

func SendMessageRequestToModel(req *telegramV1.SendMessageRequest) (sessionID, peer, text string) {
	return req.SessionId, req.Peer, req.Text
}

func MessageToProtoUpdate(msg *model.Message) *telegramV1.MessageUpdate {
	return &telegramV1.MessageUpdate{
		MessageId: msg.ID,
		From:      msg.From,
		Text:      msg.Text,
		Timestamp: msg.Timestamp,
	}
}
