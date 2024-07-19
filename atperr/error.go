package atperr

import "strings"

const (
	errorInternalServerError = "InternalServerError: Internal Server Error"

	errorInvalidDidOrPassword    = "AuthenticationRequired: Invalid identifier or password"
	errorAuthFactorTokenRequired = "AuthFactorTokenRequired: A sign in code has been sent to your email address"
	errorTokenExpired            = "ExpiredToken: Token has expired"
	errorTokenRevoked            = "Token has been revoked"

	errorProfileNotFound         = "InvalidRequest: Profile not found"
	errorInvalidActorDidOrHandle = "InvalidRequest: Error: actor must be a valid did or a handle"
	errorHandleMustBeValidHandle = "InvalidRequest: Error: handle must be a valid handle"
	errorParamMustHavePropHandle = `InvalidRequest: Error: Params must have the property "handle"`
	errorParamMustHavePropActor  = `InvalidRequest: Error: Params must have the property "actor"`
	errorAccountDeactivated      = "AccountDeactivated: Account is deactivated"
	errorInvalidFollowDid        = "Record/subject must be a valid did"

	errorRecipientNotFollowingYou = "InvalidRequest: recipient requires incoming messages to come from someone they follow"

	errorUpstreamFailure = "UpstreamFailure: Upstream Failure"
	errorUpstreamTimeout = "UpstreamTimeout: Upload timed out, please try again"

	errorInvalidRepo          = "InvalidRequest: Error: repo must be a valid did or a handle"
	errorCouldNotFindRepo     = "InvalidRequest: Could not find repo"
	errorCouldNotLocateRecord = "InvalidRequest: Could not locate record"
	errorBlobTooLarge         = "BlobTooLarge: This file is too large"
	errorCouldNotFindBlob     = "BlobNotFound: Could not find blob"
)

func IsInternalServerError(err error) bool {
	return strings.Contains(err.Error(), errorInternalServerError)
}

func IsInvalidDIDOrPasswordError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidDidOrPassword)
}

func IsAuthFactorTokenRequiredError(err error) bool {
	return strings.Contains(err.Error(), errorAuthFactorTokenRequired)
}

func IsTokenExpiredError(err error) bool {
	return strings.Contains(err.Error(), errorTokenExpired)
}

func IsTokenRevokedError(err error) bool {
	return strings.Contains(err.Error(), errorTokenRevoked)
}

func IsProfileNotFoundError(err error) bool {
	return strings.Contains(err.Error(), errorProfileNotFound)
}

func IsInvalidActorDidOrHandleError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidActorDidOrHandle)
}

func IsHandleMustBeValidHandleError(err error) bool {
	return strings.Contains(err.Error(), errorHandleMustBeValidHandle)
}

func IsParamMustHavePropHandleError(err error) bool {
	return strings.Contains(err.Error(), errorParamMustHavePropHandle)
}

func IsParamMustHavePropActorError(err error) bool {
	return strings.Contains(err.Error(), errorParamMustHavePropActor)
}

func IsAccountDeactivatedError(err error) bool {
	return strings.Contains(err.Error(), errorAccountDeactivated)
}

func IsInvalidFollowDidError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidFollowDid)
}

func IsRecipientNotFollowingYouError(err error) bool {
	return strings.Contains(err.Error(), errorRecipientNotFollowingYou)
}

func IsUpstreamFailureError(err error) bool {
	return strings.Contains(err.Error(), errorUpstreamFailure)
}

func IsUpstreamTimeoutError(err error) bool {
	return strings.Contains(err.Error(), errorUpstreamTimeout)
}

func IsInvalidRepoError(err error) bool {
	return strings.Contains(err.Error(), errorInvalidRepo)
}

func IsCouldNotFindRepoError(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotFindRepo)
}

func IsCouldNotLocateRecordError(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotLocateRecord)
}

func IsBlobTooLargeError(err error) bool {
	return strings.Contains(err.Error(), errorBlobTooLarge)
}

func IsCouldNotFindBlobError(err error) bool {
	return strings.Contains(err.Error(), errorCouldNotFindBlob)
}
