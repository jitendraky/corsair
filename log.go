package main

func InitLogging(){
	// Set up process log before anything bad happens
	switch logfile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "":
		log.SetOutput(ioutil.Discard)
	default:
		log.SetOutput(&lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    100,
			MaxAge:     14,
			MaxBackups: 10,
		})
	}
}

// mustLogFatalf wraps log.Fatalf() in a way that ensures the
// output is always printed to stderr so the user can see it
// if the user is still there, even if the process log was not
// enabled. If this process is an upgrade, however, and the user
// might not be there anymore, this just logs to the process
// log and exits.
func mustLogFatalf(format string, args ...interface{}) {
	if !corsair.IsUpgrade() {
		log.SetOutput(os.Stderr)
	}
	log.Fatalf(format, args...)
}


