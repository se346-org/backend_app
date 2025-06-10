package domain

const (
	//stream name
	STREAM_NAME_WS_MESSAGE   = "WS_MESSAGE"
	STREAM_NAME_CONVERSATION = "CONVERSATION"
	STREAM_NAME_FCM          = "FCM"

	//subject for conversation
	SUBJECT_WILDCARD_CONVERSATION  = "conversation.*"
	SUBJECT_UPDATE_LAST_MESSAGE_ID = "conversation.update_last_message_id"

	//subject for websocket
	SUBJECT_WILDCARD_MESSAGE          = "ws_message.*"
	SUBJECT_NEW_MESSAGE               = "ws_message.new"
	SUBJECT_WS_UPDATE_LAST_MESSAGE_ID = "ws_message.update_last_message_id"

	// consumer name for websocket
	CONSUMER_NAME_WS_MESSAGE_NEW                 = "ws_message_new_consumer"
	CONSUMER_NAME_WS_MESSAGE_UPDATE_LAST_MESSAGE = "ws_message_update_last_message_consumer"
	CONSUMER_NAME_SEEN_MESSAGE                   = "seen_message_consumer"
	//queue name
	QUEUE_NAME_WS_MESSAGE_UPDATE_LAST_MESSAGE = "ws_message_update_last_message_queue"
	QUEUE_NAME_SEEN_MESSAGE                   = "seen_message_queue"

	//subject for seen message
	SUBJECT_SEEN_MESSAGE = "conversation.seen_message"

	//subject for fcm
	SUBJECT_FCM_MESSAGE = "fcm.message"
	SUBJECT_WILDCARD_FCM = "fcm.*"

	//consumer name for fcm
	CONSUMER_NAME_FCM_MESSAGE = "fcm_message_consumer"

	//queue name for fcm
	QUEUE_NAME_FCM_MESSAGE = "fcm_message_queue"
)
