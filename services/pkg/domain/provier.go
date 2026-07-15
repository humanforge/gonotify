package domain

type Provider string

const (
	ProviderSendGrid  Provider = "SENDGRID"
	ProviderAmazonSES Provider = "AMAZON_SES"
	ProviderTwilio    Provider = "TWILIO"
	ProviderMSG91     Provider = "MSG91"
	ProviderFCM       Provider = "FCM"
	ProviderAPNS      Provider = "APNS"
)

/*
* FailureClassifier maps a provider error into temporary/permanent.
Each provider adapter package will have its own detailed classifier,
this is just the shared vocabulary/interface shape.
*/
type FailureClassifier interface {
	Classify(err error) FailureClass
}
