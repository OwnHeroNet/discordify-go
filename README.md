Discordify is a wrapper to execute UNIX shell commands and notify a channel in either Slack or Discord about the results. This is the `Go` language implementation.

## Installation

`go get github.com/OwnHeroNet/discordify-go`

## Configuration

Create a `.disco.yaml` file in your `$HOME` directory, or the current working directory. For now, there is only one configuration setting:

```yaml
webhook: https://your.chat.instance/api/webhooks/uid/token
```

Alternatively, you can skip the configuration file and pass `--webhook https://your.chat.instance/api/webhooks/uid/token` to the `discordify command`.

## Usage

Execute any shell command as usual wrapped in the discordify executable:

`discordify --webhook https://your.chat.instance/api/webhooks/uid/token -- make -j8`