package sendLetter

import (
	"errors"
	"fmt"
	"net/smtp"
	"to-do-list/app/configs"

	"github.com/jordan-wright/email"
)

func SendByEmail(emailRecipient string, tempCode uint, conf *configs.SendEmail) error {
	newSend := email.NewEmail()
	newSend.From = conf.SenderEmail
	newSend.To = []string{emailRecipient}
	newSend.Subject = "Verification letter for authorization on the To-Do-List service"
	newSend.Text = []byte(fmt.Sprintf("enter the following code %d in the password field to log in", tempCode))
	errSend := newSend.Send(conf.AddressHost, smtp.PlainAuth("", conf.SenderEmail, conf.Password, conf.Address))
	if errSend != nil {
		return errors.New("it was not possible to send a letter to email: " + emailRecipient)
	}
	return nil
}

//e := email.NewEmail()
//e.From = vr.Conf.EmailApi
//e.To = []string{resVerify.Email}
//e.Subject = "Verification message"
//e.Text = []byte(resVerify.Hash)
//e.HTML = []byte(fmt.Sprintf(`<a href="http://localhost:8081/verify/%s"> Кликните, чтобы подтвердить вход в аккаунт</a>`, resVerify.Hash))
//errSend := e.Send(vr.Conf.AddressHost, smtp.PlainAuth("", vr.Conf.EmailApi, vr.Conf.Password, vr.Conf.Address))
//if errSend != nil {
//responsejs.RespJs(w, NewResponseSend(errSend), http.StatusInternalServerError)
//return
//}
