package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// trapSignalsPosix captures POSIX-only signals.
func trapSignalsPosix() {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)

		for sig := range sigchan {
			switch sig {
			case syscall.SIGTERM:
				log.Println("[INFO] SIGTERM: Terminating process")
				if PidFile != "" {
					os.Remove(PidFile)
				}
        //TODO:Handle other exits, do not let the software stop
				//os.Exit(0)

			case syscall.SIGQUIT:
				log.Println("[INFO] SIGQUIT: Shutting down")
				exitCode := executeShutdownCallbacks("SIGQUIT")
				err := Stop()
				if err != nil {
					log.Printf("[ERROR] SIGQUIT stop: %v", err)
					exitCode = 1
				}
				if PidFile != "" {
					os.Remove(PidFile)
				}
				//os.Exit(exitCode)

			case syscall.SIGHUP:
				log.Println("[INFO] SIGHUP: Hanging up")
				err := Stop()
				if err != nil {
					log.Printf("[ERROR] SIGHUP stop: %v", err)
				}

			case syscall.SIGUSR1:
				log.Println("[INFO] SIGUSR1: Reloading")

				// Start with the existing Corsairfile
				instancesMu.Lock()
				if len(instances) == 0 {
					instancesMu.Unlock()
					log.Println("[ERROR] SIGUSR1: No server instances are fully running")
					continue
				}
				inst := instances[0] // we only support one instance at this time
				instancesMu.Unlock()

				updatedcorsairfile := inst.corsairfileInput
				if updatedCorsairfile == nil {
					// Hmm, did spawing process forget to close stdin? Anyhow, this is unusual.
					log.Println("[ERROR] SIGUSR1: no Corsairfile to reload (was stdin left open?)")
					continue
				}
				if loaderUsed.loader == nil {
					// This also should never happen
					log.Println("[ERROR] SIGUSR1: no Corsairfile loader with which to reload Corsairfile")
					continue
				}

				// Load the updated Corsairfile
				newCorsairfile, err := loaderUsed.loader.Load(inst.serverType)
				if err != nil {
					log.Printf("[ERROR] SIGUSR1: loading updated Corsairfile: %v", err)
					continue
				}
				if newCorsairfile != nil {
					updatedCorsairfile = newCorsairfile
				}

				// Kick off the restart; our work is done
				inst, err = inst.Restart(updatedCorsairfile)
				if err != nil {
					log.Printf("[ERROR] SIGUSR1: %v", err)
				}
			}
		}
	}()
}
