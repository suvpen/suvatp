package atperr

import "strings"

const (
	ErrorInternalServerErrorMessage = "InternalServerError: Internal Server Error"

	ErrorInvalidDIDOrPasswordMessage    = "AuthenticationRequired: Invalid identifier or password"
	ErrorAuthFactorTokenRequiredMessage = "AuthFactorTokenRequired: A sign in code has been sent to your email address"
	ErrorTokenExpiredMessage            = "ExpiredToken: Token has expired"
	ErrorTokenRevokedMessage            = "Token has been revoked"

	ErrorProfileNotFoundMessage         = "InvalidRequest: Profile not found"
	ErrorInvalidActorDidOrHandleMessage = "InvalidRequest: Error: actor must be a valid did or a handle"
	ErrorAccountDeactivatedMessage      = "AccountDeactivated: Account is deactivated"
	ErrorInvalidFollowDIDMessage        = "Record/subject must be a valid did"
	ErrorParamMustHavePropActorMessage  = `Params must have the property`

	ErrorUpstreamFailureMessage = "UpstreamFailure: Upstream Failure"
	ErrorUpstreamTimeoutMessage = "UpstreamTimeout: Upload timed out, please try again"

	ErrorInvalidRepoMessage          = "InvalidRequest: Error: repo must be a valid did or a handle"
	ErrorCouldNotFindRepoMessage     = "InvalidRequest: Could not find repo"
	ErrorCouldNotLocateRecordMessage = "InvalidRequest: Could not locate record"
	ErrorBlobTooLargeMessage         = "BlobTooLarge: This file is too large"
)

func ErrorInternalServerError(err error) bool {
	return strings.Contains(err.Error(), ErrorInternalServerErrorMessage)
}

func ErrorInvalidDIDOrPassword(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidDIDOrPasswordMessage)
}

func ErrorAuthFactorTokenRequired(err error) bool {
	return strings.Contains(err.Error(), ErrorAuthFactorTokenRequiredMessage)
}

func ErrorTokenExpired(err error) bool {
	return strings.Contains(err.Error(), ErrorTokenExpiredMessage)
}

func ErrorTokenRevoked(err error) bool {
	return strings.Contains(err.Error(), ErrorTokenRevokedMessage)
}

func ErrorProfileNotFound(err error) bool {
	return strings.Contains(err.Error(), ErrorProfileNotFoundMessage)
}

func ErrorInvalidActorDidOrHandle(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidActorDidOrHandleMessage)
}

func ErrorAccountDeactivated(err error) bool {
	return strings.Contains(err.Error(), ErrorAccountDeactivatedMessage)
}

func ErrorInvalidFollowDID(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidFollowDIDMessage)
}

func ErrorParamMustHavePropActor(err error) bool {
	return strings.Contains(err.Error(), ErrorParamMustHavePropActorMessage)
}

func ErrorUpstreamFailure(err error) bool {
	return strings.Contains(err.Error(), ErrorUpstreamFailureMessage)
}

func ErrorUpstreamTimeout(err error) bool {
	return strings.Contains(err.Error(), ErrorUpstreamTimeoutMessage)
}

func ErrorInvalidRepo(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidRepoMessage)
}

func ErrorCouldNotFindRepo(err error) bool {
	return strings.Contains(err.Error(), ErrorCouldNotFindRepoMessage)
}

func ErrorCouldNotLocateRecord(err error) bool {
	return strings.Contains(err.Error(), ErrorCouldNotLocateRecordMessage)
}

func ErrorBlobTooLarge(err error) bool {
	return strings.Contains(err.Error(), ErrorBlobTooLargeMessage)
}
