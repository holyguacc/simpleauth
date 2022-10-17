package api

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	db "simpleauth/db/sqlc"
)

type Request struct {
	to      []string
	subject string
}

func createNewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

// any object that has a create() function going to be the data holder
type Template interface {
	Create(string, string, string)
}

type Templatedata struct {
	FirstName  string
	MainText   string
	Link       string
	BottomText string
}

func (t *Templatedata) Create(fullname string, maintext string, footer string) {
	t.FirstName = fullname
	t.MainText = maintext
	t.BottomText = footer
}

// any object that implements email maker interface can be used to send email

type EmailMaker interface {
	Populate(string, string, db.User) *Request // creates To , from , subject and url
	CreateBodyTemplate()                       // creates template so later it can be made to an html template
	EmbedHtml() string                         // turns body struct to an html and returns it as a string
}

// this object will create a verification email
type verificationEmail struct {
	user db.User
	req  *Request
	data Templatedata
}

func (v *verificationEmail) Populate(verificatioCode string, clientUrl string, user db.User) *Request {
	v.data.Link = fmt.Sprintf("%s/verify/send?username=%s&code=%s", clientUrl, user.Username, verificatioCode)
	subject := "Verification email"
	v.user = user
	v.req = createNewRequest([]string{v.user.Email}, subject)
	return v.req
}
func (v *verificationEmail) CreateBodyTemplate() {
	mainText := "Please click on the link below to confirm your email address and finish the signing up process"
	footerText := "please ignore this email if it isnt for you"
	v.data.Create(v.user.FullName, mainText, footerText)
}

func (v *verificationEmail) EmbedHtml() string {
	wd, _ := os.Getwd()
	temp, err := template.ParseFiles(wd + "/layout/verify.html")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	temp.Execute(buf, v.data)
	return buf.String()
}

//////

// any object that has a create() function going to be the data holder

type ResetTemplatedata struct {
	FirstName  string
	MainText   string
	Link       string
	BottomText string
}

func (rt *ResetTemplatedata) Create(fullname string, maintext string, footer string) {
	rt.FirstName = fullname
	rt.MainText = maintext
	rt.BottomText = footer
}

type ResetEmail struct {
	user db.User
	req  *Request
	data ResetTemplatedata
}

func (r *ResetEmail) Populate(resetcode string, clientUrl string, user db.User) *Request {
	r.data.Link = resetcode
	subject := "Verification email"
	r.user = user
	r.req = createNewRequest([]string{r.user.Email}, subject)
	return r.req
}
func (r *ResetEmail) CreateBodyTemplate() {
	mainText := "Please click on the link below to Reset your Password"
	footerText := "please ignore this email if it isnt for you"
	r.data.Create(r.user.FullName, mainText, footerText)
}

func (r *ResetEmail) EmbedHtml() string {
	wd, _ := os.Getwd()
	temp, err := template.ParseFiles(wd + "/layout/reset.html")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	temp.Execute(buf, r.data)
	return buf.String()
}

func (server *Server) sendEmail(user db.User, code string, payload EmailMaker) error {
	request := payload.Populate(code, server.config.ClientUrl, user)
	payload.CreateBodyTemplate()
	body := payload.EmbedHtml()
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(request.to[0] + mime + "\n" + body)
	address := fmt.Sprintf("%s:%d", server.config.EmailServer, server.config.EmailServerPort)
	auth := smtp.PlainAuth("", "", "", server.config.EmailServer)
	if err := smtp.SendMail(address, auth, server.config.EmailSenderAddress, request.to, msg); err != nil {
		return fmt.Errorf("failed to send confirmation email %s", err)
	}
	return nil
}
