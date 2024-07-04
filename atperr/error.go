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

	ErrorRecipientNotFollowingYou = "InvalidRequest: recipient requires incoming messages to come from someone they follow"

	ErrorUpstreamFailureMessage = "UpstreamFailure: Upstream Failure"
	ErrorUpstreamTimeoutMessage = "UpstreamTimeout: Upload timed out, please try again"

	ErrorInvalidRepoMessage          = "InvalidRequest: Error: repo must be a valid did or a handle"
	ErrorCouldNotFindRepoMessage     = "InvalidRequest: Could not find repo"
	ErrorCouldNotLocateRecordMessage = "InvalidRequest: Could not locate record"
	ErrorBlobTooLargeMessage         = "BlobTooLarge: This file is too large"
)

func IsInternalServerError(err error) bool {
	return strings.Contains(err.Error(), ErrorInternalServerErrorMessage)
}

func IsInvalidDIDOrPasswordError(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidDIDOrPasswordMessage)
}

func IsAuthFactorTokenRequiredError(err error) bool {
	return strings.Contains(err.Error(), ErrorAuthFactorTokenRequiredMessage)
}

func IsTokenExpiredError(err error) bool {
	return strings.Contains(err.Error(), ErrorTokenExpiredMessage)
}

func IsTokenRevokedError(err error) bool {
	return strings.Contains(err.Error(), ErrorTokenRevokedMessage)
}

func IsProfileNotFoundError(err error) bool {
	return strings.Contains(err.Error(), ErrorProfileNotFoundMessage)
}

func IsInvalidActorDidOrHandleError(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidActorDidOrHandleMessage)
}

func IsAccountDeactivatedError(err error) bool {
	return strings.Contains(err.Error(), ErrorAccountDeactivatedMessage)
}

func IsInvalidFollowDidError(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidFollowDIDMessage)
}

func IsParamMustHavePropActorError(err error) bool {
	return strings.Contains(err.Error(), ErrorParamMustHavePropActorMessage)
}

func IsRecipientNotFollowingYouError(err error) bool {
	return strings.Contains(err.Error(), ErrorRecipientNotFollowingYou)
}

func IsUpstreamFailureError(err error) bool {
	return strings.Contains(err.Error(), ErrorUpstreamFailureMessage)
}

func IsUpstreamTimeoutError(err error) bool {
	return strings.Contains(err.Error(), ErrorUpstreamTimeoutMessage)
}

func IsInvalidRepoError(err error) bool {
	return strings.Contains(err.Error(), ErrorInvalidRepoMessage)
}

func IsCouldNotFindRepoError(err error) bool {
	return strings.Contains(err.Error(), ErrorCouldNotFindRepoMessage)
}

func IsCouldNotLocateRecordError(err error) bool {
	return strings.Contains(err.Error(), ErrorCouldNotLocateRecordMessage)
}

func IsBlobTooLargeError(err error) bool {
	return strings.Contains(err.Error(), ErrorBlobTooLargeMessage)
}
