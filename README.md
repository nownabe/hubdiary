# hubdiary

hubdiary is a CLI to write your diary and save it in a GitHub repository.

## Install

```bash
go install go.nownabe.dev/hubdiary@latest
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
  "owner": "owner",
  "repo": "repo",
  "branch": "main",
  "committer_name": "myname",
  "committer_email": "myname@users.noreply.github.com",
  "pat": "GitHub personal access token",
  "editor": "/path/to/vim"
}
```

* `owner`: Owner of repository to store diary in. Default comes from `user.name` of git config.
* `repo`: Repository to store diary in. Default is `diary`.
* `branch`: Branch of the repository to store diary in. Default is `main`.
* `committer_name`: Name of the author of commits. Default comes from `user.name` of git config.
* `committer_email`: Email of the author of commits. Default comes from `user.email` of git config.
* `pat`: GitHub Personal Access Token. If not specified in the config file, pat is given by `$GITHUB_PAT`.
* `editor`: Editor to write diary. If not specified, hubdiary uses `$EDITOR`.
