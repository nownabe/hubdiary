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
  "repository": "owner/repo",
  "user": "owner",
  "email": "owner@users.noreply.github.com",
  "pat": "xxx"
}
```

* `repository`: Repository to save diary. Default is `"diary"`.
* `user`: Name of the author of commits. Default is the user of `.gitconfig`.
* `email`: Email of the author of commits. Default is the email of `.gitconfig`.
* `pat`: GitHub Personal Access Token. If not specified in the config file, pat is given by `$GITHUB_PAT`.
