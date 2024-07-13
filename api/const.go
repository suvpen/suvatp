package api

const (
	ATPDir                = ".atp"
	ATPClientAuthJsonFile = ".atp/%s_%s_auth.json"

	DefaultATProtoEndpoint = "https://bsky.social"

	DefaultProfilesCollection = "app.bsky.actor.profile"
	DefaultPostsCollection    = "app.bsky.feed.post"
	DefaultRepostsCollection  = "app.bsky.feed.repost"
	DefaultLikeCollection     = "app.bsky.feed.like"
	DefaultGraphFollowLexicon = "app.bsky.graph.follow"
	DefaultGraphBlockLexicon  = "app.bsky.graph.block"
	DefaultLabelerService     = "app.bsky.labeler.service"
	DefaultRetries            = 1
)

const (
	FilterPostsWithReplies      = "posts_with_replies"
	FilterPostsNoReplies        = "posts_no_replies"
	FilterPostsWithMedia        = "posts_with_media"
	FilterPostsAndAuthorThreads = "posts_and_author_threads"
)
