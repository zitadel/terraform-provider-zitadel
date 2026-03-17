package instance_secret_generator

import (
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings"
)

const (
	generatorTypeVar       = "generator_type"
	lengthVar              = "length"
	expiryVar              = "expiry"
	includeLowerLettersVar = "include_lower_letters"
	includeUpperLettersVar = "include_upper_letters"
	includeDigitsVar       = "include_digits"
	includeSymbolsVar      = "include_symbols"
)

var generatorTypeMap = map[string]settings.SecretGeneratorType{
	"app_secret":             settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_APP_SECRET,
	"init_code":              settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_INIT_CODE,
	"otp_email":              settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_OTP_EMAIL,
	"otp_sms":                settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_OTP_SMS,
	"password_reset_code":    settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE,
	"passwordless_init_code": settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE,
	"verify_email_code":      settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE,
	"verify_phone_code":      settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE,
}