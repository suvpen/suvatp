package atperr

import "strings"

const (
	errorInternalServerErrorMessage = "InternalServerError: Internal Server Error"

	errorInvalidDIDOrPasswordMessage    = "AuthenticationRequired: Invalid identifier or password"
	errorAuthFactorTokenRequiredMessage = "AuthFactorTokenRequired: A sign in code has been sent to your email address"
	errorTokenExpiredMessage            = "ExpiredToken: Token has expired"
	errorTokenRevokedMessage            = "Token has been revoked"

	errorProfileNotFoundMessage         = "InvalidRequest: Profile not found"
	errorInvalidActorDidOrHandleMessage = "InvalidRequest: Error: actor must be a valid did or a handle"
	errorHandleMustBeValidHandle        = "InvalidRequest: Error: handle must be a valid handle"
	errorAccountDeactivatedMessage      = "AccountDeactivated: Account is deactivated"
	errorInvalidFollowDIDMessage        = "Record/subject must be a valid did"
	errorParamMustHavePropActorMessage  = `Params must have the property`

	errorRecipientNotFollowingYou = "InvalidRequest: recipient requires incoming messages to come from someone they follow"

	errorUpstreamFailureMessage = "UpstreamFailure: Upstream Failure"
	errorUpstreamTimeoutMessage = "UpstreamTimeout: Upload timed out, please try again"

	errorInvalidRepoMessage          = "InvalidRequest: Error: repo must be a valid did or a handle"
	errorCouldNotFindRepoMessage     = "InvalidRequest: Could not find repo"
	errorCouldNotLocateRecordMessage = "InvalidRequest: Could not locate record"
	errorBlobTooLargeMessage         = "BlobTooLarge: This file is too large"
	errorCouldNotFindBlob            = "BlobNotFound: Could not find blob"
)

func IsInternalServerError(err error) bool {
	return strings.Contains(err.Error(), errorInternalServerErrorMessage)
}

func IsInvalidDIDOrPasswordError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidDIDOrPasswordMessage)
}

func IsAuthFactorTokenRequiredError(err error) bool {
	return strings.Contains(err.Error(), errorAuthFactorTokenRequiredMessage)
}

func IsTokenExpiredError(err error) bool {
	return strings.Contains(err.Error(), errorTokenExpiredMessage)
}

func IsTokenRevokedError(err error) bool {
	return strings.Contains(err.Error(), errorTokenRevokedMessage)
}

func IsProfileNotFoundError(err error) bool {
	return strings.Contains(err.Error(), errorProfileNotFoundMessage)
}

func IsInvalidActorDidOrHandleError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidActorDidOrHandleMessage)
}

func IsHandleMustBeValidHandle(err error) bool {
	return strings.Contains(err.Error(), errorHandleMustBeValidHandle)
}

func IsAccountDeactivatedError(err error) bool {
	return strings.Contains(err.Error(), errorAccountDeactivatedMessage)
}

func IsInvalidFollowDidError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidFollowDIDMessage)
}

func IsParamMustHavePropActorError(err error) bool {
	return strings.Contains(err.Error(), errorParamMustHavePropActorMessage)
}

func IsRecipientNotFollowingYouError(err error) bool {
	return strings.Contains(err.Error(), errorRecipientNotFollowingYou)
}

func IsUpstreamFailureError(err error) bool {
	return strings.Contains(err.Error(), errorUpstreamFailureMessage)
}

func IsUpstreamTimeoutError(err error) bool {
	return strings.Contains(err.Error(), errorUpstreamTimeoutMessage)
}

func IsInvalidRepoError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidRepoMessage)
}

func IsCouldNotFindRepoError(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotFindRepoMessage)
}

func IsCouldNotLocateRecordError(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotLocateRecordMessage)
}

func IsBlobTooLargeError(err error) bool {
	return strings.Contains(err.Error(), errorBlobTooLargeMessage)
}

func IsCouldNotFindBlob(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotFindBlob)
}
