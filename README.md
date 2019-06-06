# Publy and Subly - Publish/Receive messages to/from Pub/Sub

Quick and dirty clients for Pub/Sub testing.

## Build

```
make
```

## Publy
```
NAME:
   publy - Publish messages to Pub/Sub

USAGE:
   publy [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --project value  required - project to use
   --topic value    required - topic to publish to
   --message value  required - body of the message to send
   --create         optional - create topic if it doesn't exist
   --list           optional - list topics
   --help, -h       show help
   --version, -v    print the version
```

## Subly
```
NAME:
   subly - Receive messages from Pub/Sub

USAGE:
   subly [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --project value       required - project to use
   --subscription value  required - subscription to use
   --timeout value       optional - timeout after seconds waiting for messages (default: 10)
   --create              optional - create topic if it doesn't exist
   --list                optional - list subscriptions for topic
   --topic value         required if create or list are true - topic to publish to
   --help, -h            show help
   --version, -v         print the version
```
