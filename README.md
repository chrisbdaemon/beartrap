# Beartrap v0.3
==============================

Beartrap is meant to be a portable network defense utility written entirely in
Go.  It opens "trigger" ports on the host that an attacker would connect to.
When the attacker connects and/or performs some interactions with the trigger
an alert is raised and the attacker's ip address is either blacklisted, logged,
or some other action is taken.

The idea came from listening to the PaulDotCom Security Podcast
(http://pauldotcom.com/security-weekly/) particularly episodes 203 and 204 as
the concept of honeyports is described.

**Features:**

- Beartrap can be build for most any major platform
- Beartrap is built to let one easily create or customize different trap types
  or handlers so you can decide what happens when the trap is sprung.c