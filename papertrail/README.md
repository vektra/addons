# github.com/vektra/addons/papertrail

## Input

Plain text TCP, which must be opted into in the Papertrail UI.

## Format

Syslog RFC 5424 inspired

Papertrail supports and automatically detects syslog format messages. Syslog is
documented as RFC 5424 and RFC 3164. We send the logs using the RFC 5424 style.

In Papertrail's Event viewer, the sender name and program/component name become
clickable orange and blue links to see surrounding context.

For example, the message:

```
<22>1 2014-06-18T09:56:21Z sendername programname - - - the log message
```

Is displayed like this in Papertrail's log viewer:

```
Jun 18 9:56:21 sendername programname: the log message
```

For more information see
[Configuring remote syslog from embedded or proprietary
systems](http://help.papertrailapp.com/kb/configuration/configuring-remote-syslog-from-embedded-or-proprietary-systems/).
