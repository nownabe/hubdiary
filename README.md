# hubdiary

hubdiary is a CLI to write your diary and save it in a GitHub repository.

## Install

```bash
go install go.nownabe.dev/hubdiary
```

## Usage

Set GitHub Personal Access Token to the environment variable `GITHUB_PAT`.

```bash
export GITHUB_PAT="xxx"
```

```bash
hubdiary

# or you can specify the date
hubdiary 2006-01-02
```

## Configuration

You can configure hubdiary with `$XDG_CONFIG_HOME/hubdiary/config.json`.
If you don't specify `$XDG_CONFIG_HOME`, `~/.config/hubdiary/config.json` is the default path to the config file.

This is an example.

```json
{
  "repository": "github.com/owner/repo",
  "user": "owner",
  "email" "owner@users.noreply.github.com"
}
```

* `repository`: Repository to save diary.
* `user`: Name of the author of commits.
* `email`: Email of the author of commits.
