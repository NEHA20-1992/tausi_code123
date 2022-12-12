package helper

import (
	"fmt"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

func SendEmailService(recipient, subject, htmlBodyContent string) (err error) {

	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	//Sender := "Support <support@3edge.in>"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	//Recipient = "Muralikrishnan S <muralikrishnan.s@3edge.in>"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	//Subject = "Amazon SES Test (AWS SDK for Go)"

	// The HTML body for the email.
	// HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
	// 	"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
	// 	"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	//var TextBody string = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	var CharSet string = "UTF-8"

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		//Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			config.ServerConfiguration.Amazonses.AccessKeyID,
			config.ServerConfiguration.Amazonses.SecretAccessKey,
			"")})

	if err != nil {
		displayError(err)
		panic(err)
	}

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBodyContent),
				},
				// Text: &ses.Content{
				// 	Charset: aws.String(CharSet),
				// 	Data:    aws.String(TextBody),
				// },
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(config.ServerConfiguration.Amazonses.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		displayError(err)
		return
	}

	//fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)

	return
}

func SendEmailServiceSmtp(recipient string, name string, subject string, htmlBodyContent string, result *model.User, templateName int) (err error) {

	m := gomail.NewMessage()
	m.SetHeader("From", "tausi.score@gmail.com")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBodyContent)

	// Send the email to Bob
	d := gomail.NewDialer("smtp.gmail.com", 465, "tausi.score@gmail.com", "tausiadmin")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	// // Sender data.
	// from := "tausi.score@gmail.com"
	// password := "tausiadmin"

	// // Receiver email address.
	// to := []string{
	// 	recipient,
	// }

	// // smtp server configuration.
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// // Authentication.
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// templatePath := filepath.Clean(filepath.Join(".", "template", "welcometempl.html"))

	// t, _ := template.ParseFiles(templatePath)
	// if templateName == 2 {
	// 	t, _ = template.ParseFiles("./template/forgetpasswordtempl.html")
	// }

	// var body bytes.Buffer

	// mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// body.Write([]byte(fmt.Sprintf("Subject: "+subject+" \n%s\n\n", mimeHeaders)))

	// t.Execute(&body, struct {
	// 	FirstName string
	// 	Email     string
	// 	Url       string
	// }{
	// 	FirstName: result.FirstName,
	// 	Email:     result.Email,
	// 	Url:       config.ServerConfiguration.Amazonses.PasswordResetUrl + "?email=" + result.Email + "&resetCode=" + result.ResetCode,
	// })

	// // Sending email.
	// err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	return
}

func displayError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
}
