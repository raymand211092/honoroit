package config

var defaultConfig = &Config{
	Prefix:    "!ho",
	LogLevel:  "INFO",
	CacheSize: 2000,
	DB: DB{
		DSN:     "/tmp/honoroit.db",
		Dialect: "sqlite3",
	},
	Text: Text{
		PrefixOpen: "[OPEN]",
		PrefixDone: "[DONE]",
		Greetings:  "Hello\nyour message was sent to operators. Please, keep calm and wait for answer, it takes 1-2 days.\nDon't forget that instant messenger is the same communication channel as email, so don't expect an instant response.",
		Error:      "Something is wrong.\nI already notified developers, they're fixing the issue.\n\nPlease, try again later or use any other contact method.",
		EmptyRoom:  "Customer left the room.\nConsider that request as closed.",
		Start:      "Customer has been invited to the new room. Send messages into that thread and they will be automatically forwarded.",
		Done:       "Operator marked your request as completed.\nIf you think that it's not done yet, please start another 1:1 chat with me to open a new request.",
	},
}
