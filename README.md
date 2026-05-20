## Woah! What is this?
### I'm so glad you asked! A very simple to use cross-platform audio API built completely in go!
The idea is to have a simple-to-use API that works on many platforms and with many audio backends.
The tradeoff will be reduced control for the sake of simplicity. I want a minimal public API that
can get audio playing quickly without any of the cgo stuff. Yuck!

## Wow man totally righteous! Where are you at with this?
### I know right. I'm trying to get a pulse audio client implementation working (easier said than done)
I am basing my code (a.k.a Ctrl+c Ctrl+v) on [jfreymuth's PulseAudio client implementation](https://github.com/jfreymuth/pulse) but im trying to rework and do it from
the ground up (because I hate myself and like to waste time.)

## Hey! Wait a minute are you qualified and/or motivated enough to work on something like this and actually make it usable or even complete?
### No.

## Package-level Overview
### audigo
Top-level package handles re-exporting necessary internal types and platform dependant backend selection

### audigo/internal
Defines types and interfaces utilized and implemented by the various backends

### audigo/pulse
Work-in-progress pulse audio backend

### audigo/null
A null backend for testing

### audio/cmd
Executables for testing
