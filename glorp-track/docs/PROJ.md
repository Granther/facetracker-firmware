# Captive Portal
- User is directed here to fill out wifi information
- Device is given name {name}.local. Maybe user can pick
- Save wifi information
- Watchdog program, if it doesnt get an IP (does not need internet), put in setup mode 

## Entry
- Main run in systemd, do we have an IP. Yes? Dont start hotspot and start face camera publisher
- Also, before all of that, run a diagnostic. Ie, are all expected files there? Is camera there?
- If no, Start hotspot and other neccassary services for setup. Once done, restart service. If ok, will enter facetrack mode, if not (password wrong), setup will restart

## Problems
- What if user wants to reset device? No button? Maybe button, maybe later...

